package staticman

import (
	"log"
	"net/http"
	"strings"

	"github.com/sters/staticman/conf"
)

// RegisterHandler define http request handler
func RegisterHandler(s *http.ServeMux) {
	s.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {
		if !validate(res, req) {
			log.Print("failed to validate", req.Header)
			res.WriteHeader(500)
			return
		}

		reponame := strings.TrimSpace(req.Header.Get(repoKey))
		dirname := generateDirectoryName(reponame)
		repo, err := createLocalRepositroy(dirname)
		if err != nil {
			log.Print("failed to create localrepo", err)
			res.WriteHeader(500)
			return
		}
		if err := generateNewDeploySSHKey(repo); err != nil {
			log.Print("failed to create sshkey", err)
			_ = cleanRepo(repo)
			res.WriteHeader(500)
			return
		}
		if err := configureSSHKey(repo); err != nil {
			log.Print("failed to create configure sshkey", err)
			_ = cleanRepo(repo)
			res.WriteHeader(500)
			return
		}
		if err := configureOriginRepository(repo, reponame); err != nil {
			log.Print("failed to create configure origin", err)
			_ = cleanRepo(repo)
			res.WriteHeader(500)
			return
		}

		b, err := getSSHPublicKeyContent(repo)
		if err != nil {
			log.Print("failed to get public key", err)
			_ = cleanRepo(repo)
			res.WriteHeader(500)
			return
		}

		res.WriteHeader(200)
		res.Write(b)
		return
	})

	s.HandleFunc("/pull", func(res http.ResponseWriter, req *http.Request) {
		if !validate(res, req) {
			log.Print("failed to validate", req.Header)
			res.WriteHeader(500)
			return
		}

		reponame := generateDirectoryName(
			strings.TrimSpace(req.Header.Get(repoKey)),
		)
		repo, err := loadLocalRepository(reponame)
		if err != nil {
			log.Print("failed to load repo", err)
			res.WriteHeader(500)
			return
		}

		if err := doGitPull(repo); err != nil {
			log.Print("failed to gitpull", err)
			res.WriteHeader(500)
			return
		}

		res.WriteHeader(200)
		res.Write([]byte("ok"))
		return
	})
}

const (
	validateKey = "X-STATICMAN-KEY"
	repoKey     = "X-STATICMAN-REPONAME"
)

func validate(res http.ResponseWriter, req *http.Request) bool {
	return req.Header.Get(validateKey) == conf.Variables.HTTPHeaderKey ||
		strings.TrimSpace(req.Header.Get(repoKey)) == ""
}
