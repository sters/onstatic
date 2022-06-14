package plugin

import (
	context "context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/morikuni/failure"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	pluginDir     = ".onstatic"
	pluginMapName = "EntryPoint"
	pluginNameArg = "--name"
)

var (
	ErrPluginNotHandledPath = errors.New("not found path")
	pluginEmptyMessage      = &EmptyMessage{}
)

func pluginMap() map[string]plugin.Plugin {
	return map[string]plugin.Plugin{
		pluginMapName: &onstaticPluginImpl{},
	}
}

func handshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "ONSTATIC_PLUGIN",
		MagicCookieValue: "ONSTATIC_PLUGIN",
	}
}

type API interface {
	Start(context.Context) error
	Stop(context.Context) error
	Handle(context.Context, *HandleRequest) (*HandleResponse, error)
}

type onstaticPluginClientImpl struct {
	c OnstaticPluginClient
}

func (o *onstaticPluginClientImpl) Start(ctx context.Context) error {
	_, err := o.c.Start(ctx, pluginEmptyMessage)
	return err
}

func (o *onstaticPluginClientImpl) Stop(ctx context.Context) error {
	_, err := o.c.Stop(ctx, pluginEmptyMessage)
	return err
}

func (o *onstaticPluginClientImpl) Handle(ctx context.Context, req *HandleRequest) (*HandleResponse, error) {
	return o.c.Handle(ctx, req)
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
	return &onstaticPluginClientImpl{NewOnstaticPluginClient(c)}, nil
}

type PluginClient struct {
	raw  *plugin.Client
	name string
}

func (p *PluginClient) Kill() {
	p.raw.Kill()
}

func (p *PluginClient) GetAPIClient() (API, error) {
	rpcClient, err := p.raw.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(pluginMapName)
	if err != nil {
		return nil, err
	}

	cc, ok := raw.(API)
	if !ok {
		return nil, errors.New("could not convert API interface")
	}

	return cc, nil
}

// NewClient can call from only host side. Plugin must not call this func.
func NewClient(pluginFile string) *PluginClient {
	name := getPluginName(pluginFile)

	return &PluginClient{
		name: name,
		raw: plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshakeConfig(),
			Plugins:         pluginMap(),
			Cmd:             exec.Command(pluginFile, "plugin"),
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolGRPC,
			},
		}),
	}
}

func getPluginName(pluginFile string) string {
	cmd := exec.Command(pluginFile, pluginNameArg)
	reader, _ := cmd.StdoutPipe()
	defer reader.Close()

	_ = cmd.Start()
	defer func() { _ = cmd.Wait() }()

	buf, _ := ioutil.ReadAll(reader)

	return string(buf)
}

type runningPlugin struct {
	client    *PluginClient
	apiClient API
}

type pluginList struct {
	plugins map[string]*runningPlugin
	mux     sync.RWMutex
}

func (pl *pluginList) Kill() {
	pl.mux.Lock()
	defer pl.mux.Unlock()

	wg := sync.WaitGroup{}

	for _, p := range pl.plugins {
		wg.Add(1)
		go func() {
			p.apiClient.Stop(context.Background())
			p.client.Kill()
			wg.Done()
		}()
	}

	wg.Wait()

	pl.plugins = map[string]*runningPlugin{}
}

func (pl *pluginList) Add(pluginFile string) *pluginList {
	pl.mux.Lock()
	defer pl.mux.Unlock()

	p := NewClient(pluginFile)
	if _, ok := pl.plugins[p.name]; ok {
		return pl
	}

	api, _ := p.GetAPIClient()
	api.Start(context.Background())

	pl.plugins[p.name] = &runningPlugin{
		client:    p,
		apiClient: api,
	}

	return pl
}

func (pl *pluginList) Handle(ctx context.Context, path string, body string) (string, error) {
	pl.mux.RLock()
	defer pl.mux.RUnlock()

	for _, p := range pl.plugins {
		res, err := p.apiClient.Handle(ctx, &HandleRequest{
			Path: path,
			Body: body,
		})

		st, ok := status.FromError(err)
		if !ok {
			log.Printf("next err: %+v", err)
			return "", failure.Wrap(err)
		}

		if st.Message() == ErrPluginNotHandledPath.Error() {
			log.Print("next")
			continue
		}

		if err != nil {
			log.Printf("next err: %+v", err)
			return "", failure.Wrap(err)
		}

		return res.Body, nil
	}

	return "", ErrPluginNotHandledPath
}

var actualPluginList = &pluginList{
	plugins: map[string]*runningPlugin{},
}

func LoadPlugin(pluginFile string) {
	actualPluginList.Add(pluginFile)
}

func Handle(ctx context.Context, path string, body string) (string, error) {
	return actualPluginList.Handle(ctx, path, body)
}

func Kill() {
	actualPluginList.Kill()
}

// Serve can call from only plugin side. Host must not call this func.
func Serve(server OnstaticPluginServer) {
	pluginName, err := server.Name(context.Background(), pluginEmptyMessage)
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 && os.Args[1] == pluginNameArg {
		fmt.Println(pluginName.Name)
		os.Exit(1)
	}

	p := pluginMap()
	p[pluginMapName].(*onstaticPluginImpl).realServer = server

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig(),
		Plugins:         p,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
