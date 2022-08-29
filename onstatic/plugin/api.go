package plugin

import (
	context "context"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"
	grpc "google.golang.org/grpc"
)

const (
	EntryPoint = "EntryPoint"
	NameArg    = "--name"
)

var ErrPluginNotHandledPath = errors.New("not found path")

func PluginMap() map[string]plugin.Plugin {
	return map[string]plugin.Plugin{
		EntryPoint: &onstaticPluginImpl{},
	}
}

func HandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "ONSTATIC_PLUGIN",
		MagicCookieValue: "ONSTATIC_PLUGIN",
	}
}

type onstaticPluginImpl struct {
	plugin.Plugin
	realServer OnstaticPluginServer
}

func (p *onstaticPluginImpl) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterOnstaticPluginServer(s, p.realServer)
	return nil
}

func (p *onstaticPluginImpl) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return NewOnstaticPluginClient(c), nil
}

// Serve can call from only plugin side. Host must not call this func.
func Serve(server OnstaticPluginServer) {
	pluginName, err := server.Name(context.Background(), &EmptyMessage{})
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 && os.Args[1] == NameArg {
		fmt.Println(pluginName.Name)
		os.Exit(1)
	}

	p := PluginMap()
	p[EntryPoint].(*onstaticPluginImpl).realServer = server

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig(),
		Plugins:         p,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
