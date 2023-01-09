package main

import (
	"context"
	"encoding/json"
	"net/http"

	plugin "github.com/sters/onstatic/pluginapi"
)

type echo struct {
	plugin.BasicServer
}

func (*echo) Name(context.Context, *plugin.EmptyMessage) (*plugin.NameResponse, error) {
	return &plugin.NameResponse{
		Name: "echo",
	}, nil
}

func (e *echo) Start(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	e.RegisterHandler(plugin.HTTPMethodGET, "/api/echo/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(w)
		_ = enc.Encode(r.Header)
	})

	return &plugin.EmptyMessage{}, nil
}

func main() {
	plugin.Serve(&echo{})
}
