package main

import (
	"context"

	plugin "github.com/sters/onstatic/pluginapi"
)

type example struct {
	plugin.OnstaticPluginServer
}

func (*example) Name(context.Context, *plugin.EmptyMessage) (*plugin.NameResponse, error) {
	return &plugin.NameResponse{
		Name: "example",
	}, nil
}

func (*example) Start(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	plugin.Log("plugin starting")

	return &plugin.EmptyMessage{}, nil
}

func (*example) Stop(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	plugin.Log("plugin stopping")

	return &plugin.EmptyMessage{}, nil
}

func (*example) Handle(ctx context.Context, req *plugin.HandleRequest) (*plugin.HandleResponse, error) {
	switch req.Path {
	case "/api/example":
		return &plugin.HandleResponse{
			Body: "Hello, example!",
		}, nil
	}

	return nil, plugin.ErrPluginNotHandledPath
}

func main() {
	plugin.Serve(&example{})
}
