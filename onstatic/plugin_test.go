package onstatic

// TODO: implement here

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/sters/onstatic/pluginapi"
// 	oplugin "github.com/sters/onstatic/pluginapi"
// 	"go.uber.org/zap"
// 	"gopkg.in/go-git/go-billy.v4"
// 	"gopkg.in/go-git/go-billy.v4/memfs"
// )

// type fakeHttpResponseWriter struct {
// 	header int
// 	http.ResponseWriter
// }

// func (f *fakeHttpResponseWriter) WriteHeader(s int) {
// 	f.header = s
// }

// var testAPI1HandlerFoo = func(res http.ResponseWriter, req *http.Request) {
// 	res.WriteHeader(200)
// }

// var testAPI1HandlerBar = func(res http.ResponseWriter, req *http.Request) {
// 	res.WriteHeader(404)
// }

// type testAPI1 struct{}

// func (*testAPI1) Handlers() oplugin.Handlers {
// 	return oplugin.Handlers{
// 		"/foo": testAPI1HandlerFoo,
// 		"/bar": testAPI1HandlerBar,
// 	}
// }
// func (*testAPI1) Initialize(context.Context) {}
// func (*testAPI1) Stop(context.Context)       {}

// var _ plugin.API = (*testAPI1)(nil)

// func Test_handlePlugin(t *testing.T) {
// 	ctx := context.Background()

// 	// prepare
// 	fs := map[string]billy.Filesystem{}
// 	fsNew = func(dirpath string) billy.Filesystem {
// 		if f, ok := fs[dirpath]; ok {
// 			return f
// 		}
// 		fs[dirpath] = memfs.New()
// 		return fs[dirpath]
// 	}

// 	reponame := "git@github.com:sters/onstatic.git"
// 	dirname := getHashedDirectoryName(reponame)
// 	repo, err := createLocalRepositroy(dirname)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	repoFs, err := repoToFs(repo).Chroot("../")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	err = repoFs.MkdirAll(oplugin.PluginDir, os.FileMode(0777))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	plugFs, err := repoFs.Chroot(oplugin.PluginDir)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	plug, err := plugFs.Create("testapi.so")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = plug.Write([]byte("dummy"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	loadPlugin = func(repoFs billy.Filesystem, filename string) (oplugin.EntryPoint, error) {
// 		return func(_ context.Context, l *zap.Logger) oplugin.API {
// 			return &testAPI1{}
// 		}, nil
// 	}

// 	// test
// 	want := &fakeHttpResponseWriter{}
// 	testAPI1HandlerFoo(want, nil)

// 	got := &fakeHttpResponseWriter{}
// 	f := handlePlugin(ctx, fmt.Sprintf("/%s/foo", dirname))
// 	f(got, nil)
// 	if want.header != got.header {
// 		t.Fatalf("failed to wrong handle function: want = %+v, got = %+v", want, got)
// 	}
// }

// func Test_checkLastModTime(t *testing.T) {
// 	repoName := "test_repo"
// 	filename := "test.txt"

// 	fs := memfs.New()
// 	_, err := fs.Create(filename)
// 	if err != nil {
// 		t.Fatalf("failed to create on-memory file: %+v", err)
// 	}

// 	fstat, err := fs.Stat(filename)
// 	if err != nil {
// 		t.Fatalf("failed to get stat on-memory file: %+v", err)
// 	}

// 	// no plugin loaded
// 	err = checkLastModTime(fstat, repoName)
// 	if err != nil {
// 		t.Fatalf("failed to checkLastModTime: %+v", err)
// 	}

// 	// plugin loaded but old one
// 	loadedPlugins.plugins = map[string]repoPlugins{
// 		repoName: {
// 			plugins: map[string]repoPlugin{
// 				filename: {
// 					lastModified: time.Now().Add(-1 * time.Minute),
// 					handlers:     plugin.Handlers{},
// 				},
// 			},
// 		},
// 	}
// 	err = checkLastModTime(fstat, repoName)
// 	if err != nil {
// 		t.Fatalf("failed to checkLastModTime: %+v", err)
// 	}

// 	// plugin loaded and latest one
// 	loadedPlugins.plugins = map[string]repoPlugins{
// 		repoName: {
// 			plugins: map[string]repoPlugin{
// 				filename: {
// 					lastModified: time.Now().Add(time.Minute),
// 					handlers:     plugin.Handlers{},
// 				},
// 			},
// 		},
// 	}

// 	err = checkLastModTime(fstat, repoName)
// 	if err == nil {
// 		t.Fatal("state is plugin loaded and use latest but not got error")
// 	}
// }
