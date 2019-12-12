package onstatic

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/sters/onstatic/conf"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
)

func Test_getHashedDirectoryName(t *testing.T) {
	if getHashedDirectoryName("foo") == getHashedDirectoryName("bar") {
		t.Error("failed to getHashedDirectoryName")
	}
	if getHashedDirectoryName("foo") == getHashedDirectoryName("fooo") {
		t.Error("failed to getHashedDirectoryName")
	}
}

func Test_roughFunctional(t *testing.T) {
	conf.Init()

	// shared memfs
	fs := map[string]billy.Filesystem{}
	fsNew = func(dirpath string) billy.Filesystem {
		if f, ok := fs[dirpath]; ok {
			return f
		}
		fs[dirpath] = memfs.New()
		return fs[dirpath]
	}

	t.Log("createLocalRepositroy")
	reponame := "git@github.com:sters/onstatic.git"
	dirname := getHashedDirectoryName(reponame)
	repo, err := createLocalRepositroy(dirname)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("generateNewDeploySSHKey")
	if err := generateNewDeploySSHKey(repo); err != nil {
		t.Fatal(err)
	}
	// after check privatekey and publickey
	savedpubkey := []byte{}
	{
		fs, err := repoToFs(repo).Chroot(conf.Variables.KeyDirectoryRelatedFromRepository)
		if err != nil {
			t.Fatal(err)
		}

		f, err := fs.Open(conf.Variables.SSHPubKeyFilename)
		if err != nil {
			t.Fatal(err)
		}
		savedpubkey, err = ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Log("configureSSHKey")
	if err := configureSSHKey(repo); err != nil {
		t.Fatal(err)
	}

	t.Log("configureOriginRepository")
	if err := configureOriginRepository(repo, reponame); err != nil {
		t.Fatal(err)
	}

	// Always fail because it need to set DeployKey
	// t.Log("doGitPull")
	// if err := doGitPull(repo); err != nil {
	// 	t.Fatal(err)
	// }

	t.Log("loadLocalRepository")
	repo, err = loadLocalRepository(getHashedDirectoryName(reponame))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("getSSHPublicKeyContent")
	if pubkey, err := getSSHPublicKeyContent(repo); err != nil || !reflect.DeepEqual(pubkey, savedpubkey) {
		t.Errorf("want: %s", savedpubkey)
		t.Errorf("got : %s", pubkey)
		t.Fatal(err)
	}

	t.Log("checking sshconfig")
	if cfg, err := repo.Config(); err != nil {
		t.Fatal(err)
	} else {
		cmd := cfg.Raw.Section("core").Options.Get("sshCommand")
		if want := fmt.Sprintf("ssh -i %s -F /dev/null", getSSHKeyRelatedPath()); want != cmd {
			t.Error(want)
			t.Error(cmd)
			t.Fatal("not configured core.sshCommand on git/config")
		}
	}

	t.Log("cleanRepo")
	if err := cleanRepo(repo); err != nil {
		t.Fatal(err)
	}
}
