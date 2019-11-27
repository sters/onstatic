package staticman

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sters/staticman/ssh"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

const (
	// TODO: temoporary, move to env
	salt                  = "saltsaltsalt"
	sshKeySize            = 4096
	sshKeyFilename        = "id_rsa"
	sshPubKeyFilename     = "id_rsa.pub"
	repositoriesDirecotry = "../repositories/"
)

type repository struct {
	name string
	git  git.Repository
}

func generateDirectoryName(n string) string {
	s := sha1.New()
	s.Write([]byte(salt))
	s.Write([]byte(n))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func createLocalRepositroy() (*git.Repository, error) {
	// TODO: temoporary, move to env
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Clean(filepath.Join(filepath.Dir(filename), repositoriesDirecotry))

	return git.Init(filesystem.NewStorage(osfs.New(dir), cache.NewObjectLRUDefault()), nil)
}

func generateNewDeploySSHKey() {
	ssh.GenerateKey(sshKeySize, sshKeyFilename, sshPubKeyFilename)
}

func configureSSHKey() {}

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

	return w.Pull(&git.PullOptions{RemoteName: "origin"})
}
