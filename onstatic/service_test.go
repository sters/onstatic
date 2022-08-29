package onstatic

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/sters/onstatic/conf"
	"github.com/stretchr/testify/assert"
)

func Test_getHashedDirectoryName(t *testing.T) {
	assert.NotEqual(t, getHashedDirectoryName("foo"), getHashedDirectoryName("bar"))
	assert.NotEqual(t, getHashedDirectoryName("foo"), getHashedDirectoryName("fooo"))
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
	assert.NoError(t, err)

	t.Log("generateNewDeploySSHKey")
	assert.NoError(t, generateNewDeploySSHKey(repo))

	// after check privatekey and publickey
	var savedpubkey []byte
	{
		fs, err := repoToFs(repo).Chroot(conf.Variables.KeyDirectoryRelatedFromRepository)
		assert.NoError(t, err)

		f, err := fs.Open(conf.Variables.SSHPubKeyFilename)
		assert.NoError(t, err)

		savedpubkey, err = ioutil.ReadAll(f)
		assert.NoError(t, err)
	}

	t.Log("configureSSHKey")
	assert.NoError(t, configureSSHKey(repo))

	t.Log("configureOriginRepository")
	assert.NoError(t, configureOriginRepository(repo, reponame))

	// Always fail because it need to set DeployKey
	// t.Log("doGitPull")
	// if err := doGitPull(repo); err != nil {
	// 	t.Fatal(err)
	// }

	t.Log("loadLocalRepository")
	repo, err = loadLocalRepository(getHashedDirectoryName(reponame))
	assert.NoError(t, err)

	t.Log("getSSHPublicKeyContent")
	if pubkey, err := getSSHPublicKeyContent(repo); err != nil || !reflect.DeepEqual(pubkey, savedpubkey) {
		t.Errorf("want: %s", savedpubkey)
		t.Errorf("got : %s", pubkey)
		t.Fatal(err)
	}

	t.Log("checking sshconfig")
	cfg, err := repo.Config()
	assert.NoError(t, err)

	cmd := cfg.Raw.Section("core").Options.Get("sshCommand")
	if want := fmt.Sprintf("ssh -i %s -F /dev/null", getSSHKeyRelatedPath()); want != cmd {
		t.Error(want)
		t.Error(cmd)
		t.Fatal("not configured core.sshCommand on git/config")
	}

	t.Log("cleanRepo")
	assert.NoError(t, removeRepo(repo))
}
