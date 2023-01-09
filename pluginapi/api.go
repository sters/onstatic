package pluginapi

import (
	"bytes"
	context "context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hashicorp/go-plugin"
	"github.com/morikuni/failure"
	grpc "google.golang.org/grpc"
)

const (
	EntryPoint = "EntryPoint"
	NameArg    = "--name"

	HTTPMethodGET  HTTPMethod = http.MethodGet
	HTTPMethodPOST HTTPMethod = http.MethodPost
)

var ErrPluginNotHandledPath = errors.New("not found path")

type HTTPMethod string

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

func Log(format string, a ...any) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func SignalChan() chan os.Signal {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)

	return sigCh
}

type onstaticPluginImpl struct {
	plugin.Plugin
	realServer OnstaticPluginServer
}

func (p *onstaticPluginImpl) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterOnstaticPluginServer(s, p.realServer)
	return nil
}

func (p *onstaticPluginImpl) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return NewOnstaticPluginClient(c), nil
}

// Serve can call from only plugin side. Host must not call this func.
func Serve(server OnstaticPluginServer) {
	pluginName, err := server.Name(context.Background(), &EmptyMessage{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch plugin name: %s", err)
		os.Exit(1)
	}

	// Hack for server, it's normal behavior.
	if len(os.Args) == 2 && os.Args[1] == NameArg {
		fmt.Println(pluginName.Name)
		os.Exit(1)
	}

	p := PluginMap()
	plug, ok := p[EntryPoint].(*onstaticPluginImpl)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unexpected failure when plugin loading")
		return
	}
	plug.realServer = server

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig(),
		Plugins:         p,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}

func headerMap(req *HandleRequest) map[string][]string {
	result := map[string][]string{}
	for _, v := range req.Header {
		result[strings.ToLower(v.Key)] = v.Value
	}

	return result
}

type writer struct {
	buf    *bytes.Buffer
	header int
}

var _ http.ResponseWriter = (*writer)(nil)

func (w *writer) Header() http.Header {
	return nil
}

func (w *writer) Write(b []byte) (int, error) {
	n, err := w.buf.Write(b)
	if err != nil {
		return n, failure.Wrap(err)
	}
	return n, nil
}

func (w *writer) WriteHeader(n int) {
	w.header = n
}

type BasicServer struct {
	OnstaticPluginServer

	handlers map[string]map[HTTPMethod]http.HandlerFunc
}

func (*BasicServer) Name(context.Context, *EmptyMessage) (*NameResponse, error) {
	return &NameResponse{Name: "BasicServer"}, nil
}

func (*BasicServer) Start(context.Context, *EmptyMessage) (*EmptyMessage, error) {
	return &EmptyMessage{}, nil
}

func (*BasicServer) Stop(context.Context, *EmptyMessage) (*EmptyMessage, error) {
	return &EmptyMessage{}, nil
}

func (b *BasicServer) Handle(ctx context.Context, req *HandleRequest) (*HandleResponse, error) {
	header := headerMap(req)

	m, ok := header["method"]
	mm := HTTPMethodGET
	if ok && len(m) >= 1 {
		mm = HTTPMethod(strings.ToUpper(m[0]))
	}
	path := strings.ToLower(req.Path)

	handler, ok := b.handlers[path][mm]
	if !ok {
		Log("cannot find %s, %s", path, mm)
		return nil, ErrPluginNotHandledPath
	}

	httpReq, err := http.NewRequestWithContext(ctx, string(mm), req.Path, strings.NewReader(req.Body))
	if err != nil {
		Log("%+v", err)
		return nil, ErrPluginNotHandledPath
	}
	httpReq.Header = header

	w := &writer{buf: bytes.NewBufferString("")}
	handler(w, httpReq)

	return &HandleResponse{
		Body: w.buf.String(),
	}, nil
}

func (b *BasicServer) RegisterHandler(method HTTPMethod, path string, callback http.HandlerFunc) {
	if b.handlers == nil {
		b.handlers = make(map[string]map[HTTPMethod]http.HandlerFunc)
	}

	path = strings.ToLower(path)

	if b.handlers[path] == nil {
		b.handlers[path] = make(map[HTTPMethod]http.HandlerFunc)
	}

	b.handlers[path][method] = callback
}
