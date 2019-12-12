package onstatic

import (
	"bytes"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"testing"

	"github.com/sters/onstatic/conf"
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
			// fake std logger
			logbuf := bytes.NewBuffer([]byte{})
			log.SetOutput(logbuf)
			defer func() {
				if logbuf.Len() != 0 {
					t.Logf("%s", logbuf)
				}

				if !strings.Contains(logbuf.String(), test.wantLogContents) {
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

func Test_handlePull(t *testing.T) {
	conf.Init()

	// shared memfs
	fs := map[string]billy.Filesystem{}
	fsNew = func(dirpath string) billy.Filesystem {
		if f, ok := fs[dirpath]; ok {
			log.Println("exist: " + dirpath)
			return f
		}
		log.Println("not found: " + dirpath)
		fs[dirpath] = memfs.New()
		return fs[dirpath]
	}

	reponame := "git@github.com:sters/onstatic.git"

	// setup
	createLocalRepositroy(getHashedDirectoryName(reponame))

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
			"repository does not exist",
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
			"open /.git/id_rsa: no such file or directory",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// fake std logger
			logbuf := bytes.NewBuffer([]byte{})
			log.SetOutput(logbuf)
			defer func() {
				if logbuf.Len() != 0 {
					t.Logf("%s", logbuf)
				}

				if !strings.Contains(logbuf.String(), test.wantLogContents) {
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
			f, _ := fs.Create(name)
			f.Write([]byte(body))
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
		wantContentType    string
	}{
		{
			"root not found",
			generateReq("/"),
			http.StatusNotFound,
			"",
			"",
		},
		{
			"repo path not found",
			generateReq("/foo/foo.txt"),
			http.StatusNotFound,
			"",
			"",
		},
		{
			"repo root not found",
			generateReq("/" + getHashedDirectoryName(reponame) + "/"),
			http.StatusNotFound,
			"",
			"",
		},
		{
			"avoid traversal",
			generateReq("/" + getHashedDirectoryName(reponame) + "/../../foo.txt"),
			http.StatusNotFound,
			"",
			"",
		},
		{
			"file found but cant see it",
			generateReq("/" + getHashedDirectoryName(reponame) + "/foo.bin"),
			http.StatusNotFound,
			"",
			"",
		},
		{
			"file found",
			generateReq("/" + getHashedDirectoryName(reponame) + "/foo.txt"),
			http.StatusOK,
			"Hello world",
			"text/plain",
		},
		{
			"file found with html",
			generateReq("/" + getHashedDirectoryName(reponame) + "/bar.html"),
			http.StatusOK,
			"it needs text/html",
			"text/html",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// fake std logger
			logbuf := bytes.NewBuffer([]byte{})
			log.SetOutput(logbuf)
			defer func() {
				if logbuf.Len() != 0 {
					t.Logf("%s", logbuf)
				}
			}()

			res := &fakeResponse{
				header: http.Header{},
				body:   []byte{},
			}

			handleAll(res, test.req)

			if test.wantResponseStatus != res.status {
				t.Fatalf("want = %+v, got = %+v", test.wantResponseStatus, res.status)
			}

			if got := res.Header().Get("Content-Type"); test.wantContentType != got {
				t.Fatalf("want = %s, got = %s", test.wantContentType, got)
			}

			if test.wantResponseBody != string(res.body) {
				t.Fatalf("want = %s, got = %s", test.wantResponseBody, res.body)
			}
		})
	}
}
