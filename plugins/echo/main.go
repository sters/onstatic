package main

import (
	"context"
	"strings"

	plugin "github.com/sters/onstatic/onstatic/plugin"
)

type echo struct {
	plugin.OnstaticPluginServer
}

func (*echo) Name(context.Context, *plugin.EmptyMessage) (*plugin.NameResponse, error) {
	return &plugin.NameResponse{
		Name: "echo",
	}, nil
}

func (*echo) Start(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	return &plugin.EmptyMessage{}, nil
}

func (*echo) Stop(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	return &plugin.EmptyMessage{}, nil
}

func (*echo) Handle(ctx context.Context, req *plugin.HandleRequest) (*plugin.HandleResponse, error) {
	if !strings.HasPrefix(req.Path, "/api/echo/") {
		return nil, plugin.ErrPluginNotHandledPath
	}

	r := strings.Replace(req.Path, "/api/echo/", "", 1)

	return &plugin.HandleResponse{
		Body: r,
	}, nil
}

func main() {
	plugin.Serve(&echo{})
}
