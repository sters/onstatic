package staticman

import (
	"github.com/sters/staticman/conf"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_generateDirectoryName(t *testing.T) {
	if generateDirectoryName("foo") == generateDirectoryName("bar") {
		t.Error("failed to generateDirectoryName")
	}
	if generateDirectoryName("foo") == generateDirectoryName("fooo") {
		t.Error("failed to generateDirectoryName")
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
	reponame := "git@github.com:sters/staticman.git"
	dirname := generateDirectoryName(reponame)
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
	repo, err = loadLocalRepository(generateDirectoryName(reponame))
	if err != nil {
		t.Error(err)
	}

	t.Log("cleanRepo")
	if err := cleanRepo(repo); err != nil {
		t.Error(err)
	}
}
