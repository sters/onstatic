package onstatic

import (
	"io"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"testing"

	"github.com/sters/onstatic/conf"
	"github.com/sters/onstatic/testutil"
	"go.uber.org/zap/zapcore"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
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
				Method: "GET",
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
				Method: "POST",
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
				if logbuf.FilterMessage(test.wantLogContents).Len() == 0 {
					t.Fatalf("want in logs: %s", test.wantLogContents)
				}
			}()

			res := &fakeResponse{
				header: http.Header{},
				body:   []byte{},
			}

			handleRegister(res, test.req)

			if test.wantResponseStatus != res.status {
				t.Fatalf("want = %+v, got = %+v", test.wantResponseStatus, res.status)
			}

			for k, v := range test.wantResponseHeaders {
				if res.Header().Get(k) != v {
					t.Fatalf("want = %+v, got = %+v", test.wantResponseStatus, res.status)
				}
				res.Header().Del(k)
			}
			if 0 != len(res.Header()) {
				t.Fatalf("want = no more header, got = %+v", res.Header())
			}
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
	if err != nil {
		t.Fatalf("failed to create local repository: %+v", err)
	}

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
				Method: "GET",
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
				Method: "POST",
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
				if logbuf.FilterMessage(test.wantLogContents).Len() == 0 {
					t.Fatalf("want in logs: %s", test.wantLogContents)
				}
			}()

			res := &fakeResponse{
				header: http.Header{},
				body:   []byte{},
			}

			handleUnregister(res, test.req)

			if test.wantResponseStatus != res.status {
				t.Fatalf("want = %+v, got = %+v", test.wantResponseStatus, res.status)
			}

			for k, v := range test.wantResponseHeaders {
				if res.Header().Get(k) != v {
					t.Fatalf("want = %+v, got = %+v", test.wantResponseStatus, res.status)
				}
				res.Header().Del(k)
			}
			if 0 != len(res.Header()) {
				t.Fatalf("want = no more header, got = %+v", res.Header())
			}

			if res.status == http.StatusOK {
				// check repo removed
				_, e := loadLocalRepository(getHashedDirectoryName(reponame))
				if e == nil {
					t.Fatalf("want = repo can't load, got loaded")
				}
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
	if err != nil {
		t.Fatalf("failed to create local repository: %+v", err)
	}

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
				Method: "GET",
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
				Method: "POST",
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
				Method: "POST",
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
				if logbuf.FilterMessage(test.wantLogContents).Len() == 0 {
					t.Fatalf("want in logs: %s", test.wantLogContents)
				}
			}()

			res := &fakeResponse{
				header: http.Header{},
				body:   []byte{},
			}

			handlePull(res, test.req)

			if test.wantResponseStatus != res.status {
				t.Fatalf("want = %+v, got = %+v", test.wantResponseStatus, res.status)
			}

			for k, v := range test.wantResponseHeaders {
				if res.Header().Get(k) != v {
					t.Fatalf("want = %+v, got = %+v", test.wantResponseStatus, res.status)
				}
				res.Header().Del(k)
			}
			if 0 != len(res.Header()) {
				t.Fatalf("want = no more header, got = %+v", res.Header())
			}
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
		fs := fsNew(getRepositoryDirectoryPath(getHashedDirectoryName(reponame)))
		files := map[string]string{
			"/foo.txt":   "Hello world",
			"/bar.html":  "it needs text/html",
			"/.htaccess": "can't see this file",
			"/foo.bin":   "can't see this file",
		}
		for name, body := range files {
			f, err := fs.Create(name)
			if err != nil {
				t.Fatalf("failed to create file: %+v", err)
			}

			if _, err := f.Write([]byte(body)); err != nil {
				t.Fatalf("failed to write buffer: %+v", err)
			}

			f.Close()
		}
	}

	generateReq := func(p string) *http.Request {
		return &http.Request{
			Method: "GET",
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

			if test.wantResponseStatus != res.status {
				t.Fatalf("want = %+v, got = %+v", test.wantResponseStatus, res.status)
			}

			if test.wantResponseBody != string(res.body) {
				t.Fatalf("want = %s, got = %s", test.wantResponseBody, res.body)
			}
		})
	}
}
