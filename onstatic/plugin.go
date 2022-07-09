package onstatic

import (
	context "context"
	"errors"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/morikuni/failure"
	pluginpb "github.com/sters/onstatic/onstatic/plugin"
	"google.golang.org/grpc/status"
)

const (
	pluginDir    = ".onstatic"
	pluginBinary = "main"
)

var (
	errInvalidRequest = failure.Unexpected("invalid request")
)

type PluginClient struct {
	raw  *plugin.Client
	name string
}

func (p *PluginClient) Kill() {
	p.raw.Kill()
}

func (p *PluginClient) GetAPIClient() (pluginpb.OnstaticPluginClient, error) {
	rpcClient, err := p.raw.Client()
	if err != nil {
		return nil, failure.Wrap(err)
	}

	raw, err := rpcClient.Dispense(pluginpb.EntryPoint)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	cc, ok := raw.(pluginpb.OnstaticPluginClient)
	if !ok {
		return nil, errors.New("could not convert API interface")
	}

	return cc, nil
}

// NewPluginClient can call from only host side. Plugin must not call this func.
func NewPluginClient(pluginFile string) *PluginClient {
	name := getPluginName(pluginFile)

	return &PluginClient{
		name: name,
		raw: plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: pluginpb.HandshakeConfig(),
			Plugins:         pluginpb.PluginMap(),
			Cmd:             exec.Command(pluginFile, "plugin"),
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolGRPC,
			},
		}),
	}
}

func getPluginName(pluginFile string) string {
	cmd := exec.Command(pluginFile, pluginpb.NameArg)
	reader, _ := cmd.StdoutPipe()
	defer reader.Close()

	_ = cmd.Start()
	defer func() { _ = cmd.Wait() }()

	buf, _ := ioutil.ReadAll(reader)

	return string(buf)
}

type runningPlugin struct {
	client    *PluginClient
	apiClient pluginpb.OnstaticPluginClient
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
		p := p
		wg.Add(1)
		go func() {
			_, _ = p.apiClient.Stop(context.Background(), &pluginpb.EmptyMessage{})
			p.client.Kill()
			wg.Done()
		}()
	}

	wg.Wait()

	pl.plugins = map[string]*runningPlugin{}
}

func (pl *pluginList) Add(pluginFile string) (*pluginList, error) {
	pl.mux.Lock()
	defer pl.mux.Unlock()

	if _, ok := pl.plugins[pluginFile]; ok {
		return pl, nil
	}

	p := NewPluginClient(pluginFile)
	api, err := p.GetAPIClient()
	if err != nil {
		return nil, failure.Wrap(err)
	}
	_, _ = api.Start(context.Background(), &pluginpb.EmptyMessage{})

	pl.plugins[pluginFile] = &runningPlugin{
		client:    p,
		apiClient: api,
	}

	return pl, nil
}

func (pl *pluginList) Handle(ctx context.Context, path string, body string) (string, error) {
	pl.mux.RLock()
	defer pl.mux.RUnlock()

	for _, p := range pl.plugins {
		res, err := p.apiClient.Handle(ctx, &pluginpb.HandleRequest{
			Path: path,
			Body: body,
		})

		st, ok := status.FromError(err)
		if !ok {
			log.Printf("next err: %+v", err)
			return "", failure.Wrap(err)
		}

		if st.Message() == pluginpb.ErrPluginNotHandledPath.Error() {
			log.Print("next")
			continue
		}

		if err != nil {
			log.Printf("next err: %+v", err)
			return "", failure.Wrap(err)
		}

		return res.Body, nil
	}

	return "", pluginpb.ErrPluginNotHandledPath
}

var actualPluginList = &pluginList{
	plugins: map[string]*runningPlugin{},
}

func LoadPlugin(pluginFile string) error {
	_, err := actualPluginList.Add(pluginFile)
	if err != nil {
		return failure.Wrap(err)
	}

	return nil
}

func HandlePlugin(ctx context.Context, requestPath string, body string) (string, error) {
	pathes := strings.Split(requestPath, "/")
	if len(pathes) < 2 {
		return "", errInvalidRequest
	}

	repoName := pathes[1]
	pathUnderRepo := "/" + strings.Join(pathes[2:], "/")

	err := LoadPlugin(filepath.Join(repoName, pluginDir, pluginBinary))
	if err != nil {
		return "", failure.Wrap(err)
	}

	return actualPluginList.Handle(ctx, pathUnderRepo, body)
}

func KillAllPlugin() {
	actualPluginList.Kill()
}
