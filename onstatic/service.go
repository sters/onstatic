package onstatic

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/morikuni/failure"
	"github.com/sters/onstatic/conf"
	"go.uber.org/zap"
)

const originName = "origin"

var fsNew = func(dirpath string) billy.Filesystem {
	return osfs.New(dirpath)
}

func repoToFs(r *git.Repository) billy.Filesystem {
	return r.Storer.(*filesystem.Storage).Filesystem()
}

func repoToDir(r *git.Repository) string {
	return r.Storer.(*filesystem.Storage).Filesystem().Root()
}

// func getRepoFs(repoName string) billy.Filesystem {
// 	return fsNew(filepath.Join(getRepositoriesDir(), repoName))
// }

func getRepositoriesDir() string {
	d, _ := os.Getwd()
	return filepath.Clean(filepath.Join(
		d,
		conf.Variables.RepositoriesDirectory,
	))
}

func getSSHKeyRelatedPath() string {
	return filepath.Clean(filepath.Join(
		conf.Variables.KeyDirectoryRelatedFromRepository,
		conf.Variables.SSHKeyFilename,
	))
}

// func getSSHPubKeyRelatedPath() string {
// 	return filepath.Clean(filepath.Join(
// 		conf.Variables.KeyDirectoryRelatedFromRepository,
// 		conf.Variables.SSHPubKeyFilename,
// 	))
// }

func removeRepo(repo *git.Repository) error {
	w, err := repo.Worktree()
	if err != nil {
		return failure.Wrap(err)
	}

	return failure.Wrap(util.RemoveAll(w.Filesystem, "/"))
}

func getHashedDirectoryName(n string) string {
	s := sha1.New()
	if _, err := s.Write([]byte(conf.Variables.Salt)); err != nil {
		zap.L().Error("failed to create sha1", zap.Error(err))
		return ""
	}
	if _, err := s.Write([]byte(n)); err != nil {
		zap.L().Error("failed to create sha1", zap.Error(err))
		return ""
	}
	return fmt.Sprintf("%x", s.Sum(nil))
}

func getRepositoryDirectoryPath(reponame string) string {
	return filepath.Clean(filepath.Join(getRepositoriesDir(), reponame))
}

func createLocalRepositroy(reponame string) (*git.Repository, error) {
	dir := getRepositoryDirectoryPath(reponame)
	fs := fsNew(dir)
	if err := fs.MkdirAll("/", 0755); err != nil {
		return nil, failure.Wrap(err)
	}

	gitdir, err := fs.Chroot(".git")
	if err != nil {
		return nil, failure.Wrap(err)
	}

	r, err := git.Init(
		filesystem.NewStorage(
			gitdir,
			cache.NewObjectLRUDefault(),
		),
		fs,
	)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return r, nil
}

func loadLocalRepository(reponame string) (*git.Repository, error) {
	dir := getRepositoryDirectoryPath(reponame)
	fs := fsNew(dir)
	gitdir, err := fs.Chroot(".git")
	if err != nil {
		return nil, failure.Wrap(err)
	}

	r, err := git.Open(
		filesystem.NewStorage(
			gitdir,
			cache.NewObjectLRUDefault(),
		),
		fs,
	)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return r, nil
}

func generateNewDeploySSHKey(repo *git.Repository) error {
	k, err := generateKey(
		conf.Variables.SSHKeySize,
		conf.Variables.SSHKeyFilename,
		conf.Variables.SSHPubKeyFilename,
	)
	if err != nil {
		return failure.Wrap(err)
	}

	if err := k.save(repoToFs(repo), conf.Variables.KeyDirectoryRelatedFromRepository); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func configureSSHKey(repo *git.Repository) error {
	cfg, err := repo.Config()
	if err != nil {
		return failure.Wrap(err)
	}

	cfg.Raw.Section("core").AddOption(
		"sshCommand",
		fmt.Sprintf(
			"ssh -i %s -F /dev/null",
			getSSHKeyRelatedPath(),
		),
	)

	if err := repo.Storer.SetConfig(cfg); err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func getSSHPublicKeyContent(repo *git.Repository) ([]byte, error) {
	fs, err := repoToFs(repo).Chroot(conf.Variables.KeyDirectoryRelatedFromRepository)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to chroot"))
	}

	f, err := fs.Open(conf.Variables.SSHPubKeyFilename)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to open"))
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, failure.Wrap(err, failure.Message("failed to readall"))
	}
	return b, nil
}

func configureOriginRepository(repo *git.Repository, originURL string) error {
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: originName,
		URLs: []string{originURL},
	})

	if err != nil {
		return failure.Wrap(err)
	}
	return nil
}

func doGitPull(repo *git.Repository, branchName string) error {
	w, err := repo.Worktree()
	if err != nil {
		return failure.Wrap(err)
	}

	auth, err := ssh.NewPublicKeysFromFile(
		"git",
		filepath.Join(repoToDir(repo), getSSHKeyRelatedPath()),
		"", // TODO: passphrease
	)
	if err != nil {
		return failure.Wrap(err)
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: originName,
		Auth:       auth,
		Force:      true,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return failure.Wrap(err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewRemoteReferenceName(originName, branchName),
		Force:  true,
	})
	if err != nil {
		return failure.Wrap(err)
	}

	return nil
}
