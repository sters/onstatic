package onstatic

import (
	"net/http"
	"strings"

	"github.com/sters/onstatic/conf"
	"go.uber.org/zap"
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
		zap.L().Error("failed to validate", zap.Any("reqHeader", req.Header))
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	reponame := strings.TrimSpace(req.Header.Get(repoKey))
	dirname := getHashedDirectoryName(reponame)
	repo, err := createLocalRepositroy(dirname)
	if err != nil {
		zap.L().Error("failed to create localrepo", zap.Error(err))
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if err := generateNewDeploySSHKey(repo); err != nil {
		zap.L().Error("failed to create sshkey", zap.Error(err))
		_ = removeRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if err := configureSSHKey(repo); err != nil {
		zap.L().Error("failed to create configure sshkey", zap.Error(err))
		_ = removeRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if err := configureOriginRepository(repo, reponame); err != nil {
		zap.L().Error("failed to create configure origin", zap.Error(err))
		_ = removeRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	b, err := getSSHPublicKeyContent(repo)
	if err != nil {
		zap.L().Error("failed to get public key", zap.Error(err))
		_ = removeRepo(repo)
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	zap.L().Info("register success", zap.String("reponame", reponame))
	res.WriteHeader(http.StatusOK)
	res.Write(b)
	return
}

func handlePull(res http.ResponseWriter, req *http.Request) {
	if !validate(res, req) {
		zap.L().Error("failed to validate", zap.Any("reqHeader", req.Header))
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	reponame := getHashedDirectoryName(
		strings.TrimSpace(req.Header.Get(repoKey)),
	)
	repo, err := loadLocalRepository(reponame)
	if err != nil {
		zap.L().Error("failed to load repo", zap.Error(err))
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := doGitPull(repo); err != nil {
		zap.L().Error("failed to gitpull", zap.Error(err))
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	zap.L().Info("pull success", zap.String("reponame", reponame))
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(reponame))
	return
}

func handleUnregister(res http.ResponseWriter, req *http.Request) {
	if !validate(res, req) {
		zap.L().Error("failed to validate", zap.Any("reqHeader", req.Header))
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	reponame := strings.TrimSpace(req.Header.Get(repoKey))
	dirname := getHashedDirectoryName(reponame)

	repo, err := loadLocalRepository(dirname)
	if err != nil {
		zap.L().Error("failed to create localrepo", zap.Error(err))
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := removeRepo(repo); err != nil {
		zap.L().Error("failed to clean repo", zap.Error(err))
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	zap.L().Info("unregister success", zap.String("reponame", reponame))
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
	if len(pathes) <= 3 && strings.TrimSpace(pathes[len(pathes)-1]) == "" {
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
