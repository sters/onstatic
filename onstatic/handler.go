package onstatic

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sters/onstatic/conf"
)

// RegisterHandler define http request handler
func RegisterHandler(s *http.ServeMux) {
	s.HandleFunc("/register", handleRegister)
	s.HandleFunc("/pull", handlePull)
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
		_ = cleanRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if err := configureSSHKey(repo); err != nil {
		log.Print("failed to create configure sshkey: ", err)
		_ = cleanRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if err := configureOriginRepository(repo, reponame); err != nil {
		log.Print("failed to create configure origin: ", err)
		_ = cleanRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	b, err := getSSHPublicKeyContent(repo)
	if err != nil {
		log.Print("failed to get public key: ", err)
		_ = cleanRepo(repo)
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
	res.Write([]byte("ok"))
	return
}

// handleAll onstatic managing contents
func handleAll(res http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, "/.") {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	cleanedPath := path.Clean(req.URL.Path)
	pathes := strings.Split(cleanedPath, "/")
	if len(pathes) == 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	requestFilePath := strings.Replace(cleanedPath, "/"+pathes[0], "", 1)
	if len(pathes) == 1 {
		requestFilePath = "index.html"
	}

	fullpath := filepath.Clean(filepath.Join(getRepositoryDirectoryPath(pathes[0]), requestFilePath))
	if strings.Contains(fullpath, "/.") || !strings.HasPrefix(fullpath, getRepositoriesDir()) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := os.Open(fullpath)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusOK)
	if _, err := io.Copy(res, f); err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
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
