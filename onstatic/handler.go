package onstatic

import (
	"io"
	"log"
	"net/http"
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
	if cleanedPath[0] == '/' {
		cleanedPath = cleanedPath[1:]
	}

	pathes := strings.Split(cleanedPath, "/")
	if len(pathes) <= 1 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if hasIgnoreContents(cleanedPath) || hasIgnoreSuffix(cleanedPath) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	requestFilePath := strings.Replace(cleanedPath, pathes[0], "", 1)
	fs := fsNew(getRepositoryDirectoryPath(pathes[0]))
	if s, err := fs.Stat(requestFilePath); err != nil || s.IsDir() {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := fs.Open(requestFilePath)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()

	res.Header().Set("Content-Type", guessContentType(requestFilePath))
	res.WriteHeader(http.StatusOK)
	if _, err := io.Copy(res, f); err != nil {
		log.Println(err)
		return
	}
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

func guessContentType(path string) string {
	// useful one https://github.com/nginx/nginx/blob/master/conf/mime.types
	switch filepath.Ext(path) {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".gif":
		return "image/gif"
	case ".jpeg", ".jpg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".svg", ".svgz":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	}
	return "text/plain"
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
