package onstatic

import (
	"github.com/sters/onstatic/conf"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_getHashedDirectoryName(t *testing.T) {
	if getHashedDirectoryName("foo") == getHashedDirectoryName("bar") {
		t.Error("failed to getHashedDirectoryName")
	}
	if getHashedDirectoryName("foo") == getHashedDirectoryName("fooo") {
		t.Error("failed to getHashedDirectoryName")
	}
}

func TestConfigureRepository(t *testing.T) {
	conf.Init()

	filepath.Walk(getRepositoriesDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil || path == getRepositoriesDir() || strings.Contains(path, ".gitkeep") {
			return nil
		}
		t.Logf("Delete files: %s", path)
		_ = os.RemoveAll(path)
		return nil
	})

	t.Log("createLocalRepositroy")
	reponame := "git@github.com:sters/onstatic.git"
	dirname := getHashedDirectoryName(reponame)
	repo, err := createLocalRepositroy(dirname)
	if err != nil {
		t.Error(err)
	}

	t.Log("generateNewDeploySSHKey")
	if err := generateNewDeploySSHKey(repo); err != nil {
		t.Error(err)
	}

	t.Log("configureSSHKey")
	if err := configureSSHKey(repo); err != nil {
		t.Error(err)
	}

	t.Log("configureOriginRepository")
	if err := configureOriginRepository(repo, reponame); err != nil {
		t.Error(err)
	}

	// Always fail because it need to set DeployKey
	// t.Log("doGitPull")
	// if err := doGitPull(repo); err != nil {
	// 	t.Error(err)
	// }

	t.Log("loadLocalRepository")
	repo, err = loadLocalRepository(getHashedDirectoryName(reponame))
	if err != nil {
		t.Error(err)
	}

	t.Log("cleanRepo")
	if err := cleanRepo(repo); err != nil {
		t.Error(err)
	}
}
