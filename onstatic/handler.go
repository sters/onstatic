package onstatic

import (
	"log"
	"net/http"
	"strings"

	"github.com/sters/onstatic/conf"
)

var fileserver http.Handler

// RegisterHandler define http request handler
func RegisterHandler(s *http.ServeMux) {
	fileserver = http.FileServer(http.Dir(getRepositoriesDir()))

	s.HandleFunc("/register", handleRegister)
	s.HandleFunc("/pull", handlePull)
	s.HandleFunc("/unregister", handleUnregister)
	s.HandleFunc("/", handleAll)
}

func handleRegister(res http.ResponseWriter, req *http.Request) {
	if !validate(res, req) {
		log.Print("failed to validate: ", req.Header)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	reponame := strings.TrimSpace(req.Header.Get(repoKey))
	dirname := getHashedDirectoryName(reponame)
	repo, err := createLocalRepositroy(dirname)
	if err != nil {
		log.Print("failed to create localrepo: ", err)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if err := generateNewDeploySSHKey(repo); err != nil {
		log.Print("failed to create sshkey: ", err)
		_ = removeRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if err := configureSSHKey(repo); err != nil {
		log.Print("failed to create configure sshkey: ", err)
		_ = removeRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if err := configureOriginRepository(repo, reponame); err != nil {
		log.Print("failed to create configure origin: ", err)
		_ = removeRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	b, err := getSSHPublicKeyContent(repo)
	if err != nil {
		log.Print("failed to get public key: ", err)
		_ = removeRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	log.Print("register success: ", reponame)
	res.WriteHeader(http.StatusOK)
	res.Write(b)
	return
}

func handlePull(res http.ResponseWriter, req *http.Request) {
	if !validate(res, req) {
		log.Print("failed to validate: ", req.Header)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	reponame := getHashedDirectoryName(
		strings.TrimSpace(req.Header.Get(repoKey)),
	)
	repo, err := loadLocalRepository(reponame)
	if err != nil {
		log.Print("failed to load repo: ", err)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := doGitPull(repo); err != nil {
		log.Print("failed to gitpull: ", err)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	log.Print("pull success: ", reponame)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(reponame))
	return
}

func handleUnregister(res http.ResponseWriter, req *http.Request) {
	if !validate(res, req) {
		log.Print("failed to validate: ", req.Header)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	reponame := strings.TrimSpace(req.Header.Get(repoKey))
	dirname := getHashedDirectoryName(reponame)

	repo, err := loadLocalRepository(dirname)
	if err != nil {
		log.Print("failed to create localrepo: ", err)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := removeRepo(repo); err != nil {
		log.Print("failed to clean repo: ", err)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	log.Print("unregister success: ", reponame)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("ok"))
	return
}

// handleAll onstatic managing contents
func handleAll(res http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, "/.") {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	pathes := strings.Split(req.URL.Path, "/")
	if len(pathes) <= 2 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if strings.TrimSpace(pathes[len(pathes)-1]) == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if hasIgnoreContents(req.URL.Path) || hasIgnoreSuffix(req.URL.Path) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	fileserver.ServeHTTP(res, req)
}

func hasIgnoreContents(p string) bool {
	var ignoreContains = []string{
		"/.", "/internal", "/bin/",
	}
	for _, c := range ignoreContains {
		if strings.Contains(p, c) {
			return true
		}
	}
	return false
}

func hasIgnoreSuffix(p string) bool {
	var ignoreSuffix = []string{
		"/LICENSE", "/Makefile", "/README.md", "/README", "/id_rsa",
		".bin", ".exe", ".dll",
		".zip", ".gz", ".tar", ".db",
		".json", ".conf",
	}
	for _, s := range ignoreSuffix {
		if strings.HasSuffix(p, s) {
			return true
		}
	}

	return false
}

const (
	validateKey = "X-ONSTATIC-KEY"
	repoKey     = "X-ONSTATIC-REPONAME"
)

func validate(res http.ResponseWriter, req *http.Request) bool {
	return req.Header.Get(validateKey) == conf.Variables.HTTPHeaderKey &&
		strings.TrimSpace(req.Header.Get(repoKey)) != "" &&
		req.Method == http.MethodPost
}
