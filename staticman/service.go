package staticman

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"

	"gopkg.in/src-d/go-git.v4"
)

const (
	// temoporary, move to env
	salt              = "saltsaltsalt"
	sshKeyFilename    = "id_rsa"
	sshPubKeyFilename = "id_rsa.pub"
)

type repository struct {
	name string
	git  git.Repository
}

func generateDirectoryName() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	s := sha1.New()
	s.Write(bytes)

	return fmt.Sprintf("%x", s.Sum(nil))
}

func createLocalRepositroy() git.Repository {
	return git.Repository{}
}

func generateNewDeploySSHKey() string {
	return ""
}

func configureSSHKey() {}

func configureOriginRepository() {}

func doGitPull() {}
