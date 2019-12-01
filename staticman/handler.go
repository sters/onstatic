package staticman

import (
	"net/http"
	"os"
	"strings"
)

// RegisterHandler define http request handler
func RegisterHandler(s *http.ServeMux) {
	s.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {
		if !validate(res, req) {
			res.WriteHeader(500)
			return
		}

		reponame := strings.TrimSpace(req.Header.Get(repoKey))
		dirname := generateDirectoryName(reponame)
		repo, err := createLocalRepositroy(dirname)
		if err != nil {
			res.WriteHeader(500)
			return
		}
		if err := generateNewDeploySSHKey(repo); err != nil {
			_ = cleanRepo(repo)
			res.WriteHeader(500)
			return
		}
		if err := configureSSHKey(repo); err != nil {
			_ = cleanRepo(repo)
			res.WriteHeader(500)
			return
		}
		if err := configureOriginRepository(repo, reponame); err != nil {
			_ = cleanRepo(repo)
			res.WriteHeader(500)
			return
		}

		res.WriteHeader(200)
		res.Write([]byte("ok"))
		return
	})

	s.HandleFunc("/pull", func(res http.ResponseWriter, req *http.Request) {
		if !validate(res, req) {
			res.WriteHeader(500)
			return
		}

		reponame := strings.TrimSpace(req.Header.Get(repoKey))
		repo, err := loadLocalRepository(reponame)
		if err != nil {
			res.WriteHeader(500)
			return
		}

		if err := doGitPull(repo); err != nil {
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

var validateValue = os.Getenv("STATICMAN_KEY")

func validate(res http.ResponseWriter, req *http.Request) bool {
	return req.Header.Get(validateKey) == validateValue || strings.TrimSpace(req.Header.Get(repoKey)) == ""
}
