package staticman

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sters/staticman/ssh"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

const (
	// TODO: temoporary, move to env
	salt                              = "saltsaltsalt"
	sshKeySize                        = 4096
	sshKeyFilename                    = "id_rsa"
	sshPubKeyFilename                 = "id_rsa.pub"
	repositoriesDirectory             = "../repositories/"
	keyDirectoryRelatedFromRepository = "."
)

func repoToDir(r *git.Repository) string {
	return r.Storer.(*filesystem.Storage).Filesystem().Root()
}

func getRepositoriesDir() string {
	// TODO: temoporary, move to env
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(filename), repositoriesDirectory))
}

func getSSHKeyRelatedPath() string {
	return filepath.Clean(filepath.Join(keyDirectoryRelatedFromRepository, sshKeyFilename))
}
func getSSHPubKeyRelatedPath() string {
	return filepath.Clean(filepath.Join(keyDirectoryRelatedFromRepository, sshPubKeyFilename))
}

func generateDirectoryName(n string) string {
	s := sha1.New()
	s.Write([]byte(salt))
	s.Write([]byte(n))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func createLocalRepositroy(reponame string) (*git.Repository, error) {
	dir := filepath.Clean(filepath.Join(getRepositoriesDir(), reponame))
	if err := os.Mkdir(dir, 0755); err != nil {
		return nil, err
	}

	return git.Init(
		filesystem.NewStorage(
			osfs.New(filepath.Join(dir, ".git")),
			cache.NewObjectLRUDefault(),
		),
		osfs.New(filepath.Join(dir)),
	)
}

func loadLocalRepository(reponame string) (*git.Repository, error) {
	dir := filepath.Clean(filepath.Join(getRepositoriesDir(), reponame))
	return git.Open(
		filesystem.NewStorage(
			osfs.New(filepath.Join(dir, ".git")),
			cache.NewObjectLRUDefault(),
		),
		osfs.New(filepath.Join(dir)),
	)
}

func generateNewDeploySSHKey(repo *git.Repository) error {
	key, err := ssh.GenerateKey(sshKeySize, sshKeyFilename, sshPubKeyFilename)
	if err != nil {
		return err
	}

	dir := filepath.Join(repoToDir(repo), keyDirectoryRelatedFromRepository)
	if err := key.Save(dir); err != nil {
		return err
	}

	return nil
}

func configureSSHKey(repo *git.Repository) error {
	cfg, err := repo.Config()
	if err != nil {
		return err
	}

	cfg.Raw.Section("core").AddOption(
		"sshCommand",
		fmt.Sprintf(
			"ssh -i %s -F /dev/null",
			getSSHKeyRelatedPath(),
		),
	)

	return nil
}

func configureOriginRepository(repo *git.Repository, originURL string) error {
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{originURL},
	})
	if err != nil {
		return err
	}

	return nil
}

func doGitPull(repo *git.Repository) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	auth, err := gitssh.NewPublicKeysFromFile(
		"git",
		filepath.Join(repoToDir(repo), getSSHKeyRelatedPath()),
		"",
	)
	if err != nil {
		return err
	}

	return w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
}
