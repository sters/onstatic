package onstatic

import (
	"io"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/sters/onstatic/conf"
	"github.com/sters/onstatic/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

type fakeResponse struct {
	header http.Header
	body   []byte
	status int
}

func (f *fakeResponse) Header() http.Header {
	return f.header
}

func (f *fakeResponse) Write(b []byte) (int, error) {
	f.body = b
	return len(b), nil
}

func (f *fakeResponse) WriteHeader(c int) {
	f.status = c
}

type fakeFileserver struct{}

func (fake *fakeFileserver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	pathes := strings.Split(req.URL.Path, "/")
	requestFilePath := strings.Replace(req.URL.Path, "/"+pathes[1], "", 1)
	fs := fsNew(getRepositoryDirectoryPath(pathes[1]))
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

	res.WriteHeader(http.StatusOK)
	if _, err := io.Copy(res, f); err != nil {
		log.Println(err)
		return
	}
}

func Test_handleRegister(t *testing.T) {
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

	reponame := "git@github.com:sters/onstatic.git"

	tests := []struct {
		name                string
		req                 *http.Request
		wantResponseHeaders map[string]string
		wantResponseStatus  int
		wantLogContents     string
	}{
		{
			"failed to validate",
			&http.Request{
				Header: http.Header{},
				Method: http.MethodGet,
			},
			nil,
			http.StatusServiceUnavailable,
			"failed to validate",
		},
		{
			"success",
			&http.Request{
				Header: http.Header{
					textproto.CanonicalMIMEHeaderKey(validateKey): []string{conf.Variables.HTTPHeaderKey},
					textproto.CanonicalMIMEHeaderKey(repoKey):     []string{reponame},
				},
				Method: http.MethodPost,
			},
			nil,
			http.StatusOK,
			"register success",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logbuf, _ := testutil.NewLogObserver(t, zapcore.DebugLevel)
			defer func() {
				assert.NotEqual(t, 0, logbuf.FilterMessage(test.wantLogContents).Len(), "want in logs: %s", test.wantLogContents)
			}()

			res := &fakeResponse{
				header: http.Header{},
				body:   []byte{},
			}

			handleRegister(res, test.req)

			assert.Equal(t, test.wantResponseStatus, res.status)

			for k, v := range test.wantResponseHeaders {
				assert.Equal(t, v, res.Header().Get(k))
				res.Header().Del(k)
			}

			assert.Len(t, res.Header(), 0)
		})
	}
}

func Test_handleUnregister(t *testing.T) {
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

	reponame := "git@github.com:sters/onstatic.git"

	// setup
	_, err := createLocalRepositroy(getHashedDirectoryName(reponame))
	assert.NoError(t, err)

	tests := []struct {
		name                string
		req                 *http.Request
		wantResponseHeaders map[string]string
		wantResponseStatus  int
		wantLogContents     string
	}{
		{
			"failed to validate",
			&http.Request{
				Header: http.Header{},
				Method: http.MethodGet,
			},
			nil,
			http.StatusServiceUnavailable,
			"failed to validate",
		},
		{
			"success",
			&http.Request{
				Header: http.Header{
					textproto.CanonicalMIMEHeaderKey(validateKey): []string{conf.Variables.HTTPHeaderKey},
					textproto.CanonicalMIMEHeaderKey(repoKey):     []string{reponame},
				},
				Method: http.MethodPost,
			},
			nil,
			http.StatusOK,
			"unregister success",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logbuf, _ := testutil.NewLogObserver(t, zapcore.DebugLevel)
			defer func() {
				assert.NotEqual(t, 0, logbuf.FilterMessage(test.wantLogContents).Len(), "want in logs: %s", test.wantLogContents)
			}()

			res := &fakeResponse{
				header: http.Header{},
				body:   []byte{},
			}

			handleUnregister(res, test.req)

			assert.Equal(t, test.wantResponseStatus, res.status)

			for k, v := range test.wantResponseHeaders {
				assert.Equal(t, v, res.Header().Get(k))
				res.Header().Del(k)
			}

			assert.Len(t, res.Header(), 0)

			if res.status == http.StatusOK {
				// check repo removed
				_, e := loadLocalRepository(getHashedDirectoryName(reponame))
				assert.Error(t, e)
			}
		})
	}
}

func Test_handlePull(t *testing.T) {
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

	reponame := "git@github.com:sters/onstatic.git"

	// setup
	_, err := createLocalRepositroy(getHashedDirectoryName(reponame))
	assert.NoError(t, err)

	tests := []struct {
		name                string
		req                 *http.Request
		wantResponseHeaders map[string]string
		wantResponseStatus  int
		wantLogContents     string
	}{
		{
			"failed to validate",
			&http.Request{
				Header: http.Header{},
				Method: http.MethodGet,
			},
			nil,
			http.StatusServiceUnavailable,
			"failed to validate",
		},
		{
			"repo not found",
			&http.Request{
				Header: http.Header{
					textproto.CanonicalMIMEHeaderKey(validateKey): []string{conf.Variables.HTTPHeaderKey},
					textproto.CanonicalMIMEHeaderKey(repoKey):     []string{reponame + "foo"},
				},
				Method: http.MethodPost,
			},
			nil,
			http.StatusServiceUnavailable,
			"failed to load repo",
		},
		{
			"failed pull",
			&http.Request{
				Header: http.Header{
					textproto.CanonicalMIMEHeaderKey(validateKey): []string{conf.Variables.HTTPHeaderKey},
					textproto.CanonicalMIMEHeaderKey(repoKey):     []string{reponame},
				},
				Method: http.MethodPost,
			},
			nil,
			http.StatusServiceUnavailable,
			"failed to gitpull",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logbuf, _ := testutil.NewLogObserver(t, zapcore.DebugLevel)
			defer func() {
				assert.NotEqual(t, 0, logbuf.FilterMessage(test.wantLogContents).Len(), "want in logs: %s", test.wantLogContents)
			}()

			res := &fakeResponse{
				header: http.Header{},
				body:   []byte{},
			}

			handlePull(res, test.req)

			assert.Equal(t, test.wantResponseStatus, res.status)

			for k, v := range test.wantResponseHeaders {
				assert.Equal(t, v, res.Header().Get(k))
				res.Header().Del(k)
			}

			assert.Len(t, res.Header(), 0)
		})
	}
}

func Test_handleAll(t *testing.T) {
	conf.Init()
	fileserver = &fakeFileserver{}

	// shared memfs
	fs := map[string]billy.Filesystem{}
	fsNew = func(dirpath string) billy.Filesystem {
		if f, ok := fs[dirpath]; ok {
			return f
		}
		fs[dirpath] = memfs.New()
		return fs[dirpath]
	}

	reponame := "git@github.com:sters/onstatic.git"

	// setup
	{
		rootFs := fsNew(getRepositoriesDir())
		rootFs.Create(getHashedDirectoryName(reponame))

		fs := fsNew(getRepositoryDirectoryPath(getHashedDirectoryName(reponame)))
		files := map[string]string{
			"/foo.txt":   "Hello world",
			"/bar.html":  "it needs text/html",
			"/.htaccess": "can't see this file",
			"/foo.bin":   "can't see this file",
		}
		for name, body := range files {
			f, err := fs.Create(name)
			assert.NoError(t, err)

			_, err = f.Write([]byte(body))
			assert.NoError(t, err)

			f.Close()
		}
	}

	generateReq := func(p string) *http.Request {
		return &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Path: p,
			},
		}
	}

	tests := []struct {
		name               string
		req                *http.Request
		wantResponseStatus int
		wantResponseBody   string
	}{
		{
			"root not found",
			generateReq("/"),
			http.StatusNotFound,
			"",
		},
		{
			"repo path not found",
			generateReq("/foo/foo.txt"),
			http.StatusNotFound,
			"",
		},
		{
			"repo root not found",
			generateReq("/" + getHashedDirectoryName(reponame) + "/"),
			http.StatusNotFound,
			"",
		},
		{
			"avoid traversal",
			generateReq("/" + getHashedDirectoryName(reponame) + "/../../foo.txt"),
			http.StatusNotFound,
			"",
		},
		{
			"file found but cant see it",
			generateReq("/" + getHashedDirectoryName(reponame) + "/foo.bin"),
			http.StatusNotFound,
			"",
		},
		{
			"file found",
			generateReq("/" + getHashedDirectoryName(reponame) + "/foo.txt"),
			http.StatusOK,
			"Hello world",
		},
		{
			"file found with html",
			generateReq("/" + getHashedDirectoryName(reponame) + "/bar.html"),
			http.StatusOK,
			"it needs text/html",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := &fakeResponse{
				header: http.Header{},
				body:   []byte{},
			}

			handleAll(res, test.req)

			assert.Equal(t, test.wantResponseStatus, res.status)
			assert.Equal(t, test.wantResponseBody, string(res.body))
		})
	}
}
