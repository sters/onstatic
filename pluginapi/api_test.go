package pluginapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicServer_RegisterHandler(t *testing.T) {
	t.Parallel()

	server := &BasicServer{}

	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server.RegisterHandler(HTTPMethodGET, "/ABC", f)

	w := &writer{}
	server.handlers["/abc"][HTTPMethodGET](w, nil)
	assert.Equal(t, http.StatusOK, w.header)
}

func TestBasicServer_Handle_unmatch_path(t *testing.T) {
	t.Parallel()

	server := &BasicServer{}

	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server.RegisterHandler(HTTPMethodGET, "/ABC", f)

	_, err := server.Handle(context.Background(), &HandleRequest{
		Path: "/",
		Header: []*Header{
			{Key: "Method", Value: []string{http.MethodGet}},
		},
	})
	assert.Equal(t, err, ErrPluginNotHandledPath)
}

func TestBasicServer_Handle_unmatch_method(t *testing.T) {
	t.Parallel()

	server := &BasicServer{}

	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server.RegisterHandler(HTTPMethodGET, "/ABC", f)

	_, err := server.Handle(context.Background(), &HandleRequest{
		Path: "/abc",
		Header: []*Header{
			{Key: "Method", Value: []string{http.MethodPost}},
		},
	})
	assert.Equal(t, err, ErrPluginNotHandledPath)
}

func TestBasicServer_Handle(t *testing.T) {
	t.Parallel()

	server := &BasicServer{}

	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	})
	server.RegisterHandler(HTTPMethodGET, "/ABC", f)

	res, err := server.Handle(context.Background(), &HandleRequest{
		Path: "/abc",
		Header: []*Header{
			{Key: "Method", Value: []string{http.MethodGet}},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "test", res.Body)
}
