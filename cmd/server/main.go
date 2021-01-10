package main

import (
	"context"
	"net/http/cgi"
	"os"
	"os/signal"
	"syscall"

	"github.com/sters/onstatic/conf"
	"github.com/sters/onstatic/http"
	"github.com/sters/onstatic/onstatic"
	"go.uber.org/zap"
)

func main() {
	conf.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := http.NewServer(conf.Variables.HTTPPort)
	if err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
	onstatic.RegisterHandler(server.Mux)

	if conf.Variables.CGIMode {
		runCGIServerMode(ctx, server)
		return
	}

	runHTTPServerMode(ctx, server)
}

func runHTTPServerMode(ctx context.Context, server *http.Server) {
	zap.L().Info("http server starting", zap.String("HTTPPort", conf.Variables.HTTPPort))
	go func() { _ = server.Run() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	_ = server.Close()
}

func runCGIServerMode(ctx context.Context, server *http.Server) {
	if err := cgi.Serve(server.Mux); err != nil {
		zap.L().Fatal("failed to cgi.Serve", zap.Error(err))
	}

	server.Close()
}
