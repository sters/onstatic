package main

import (
	"context"
	"log"
	"net/http/cgi"
	"os"
	"os/signal"
	"syscall"

	"github.com/sters/onstatic/conf"
	"github.com/sters/onstatic/http"
	"github.com/sters/onstatic/onstatic"
)

func main() {
	conf.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := http.NewServer(conf.Variables.HTTPPort)
	if err != nil {
		log.Fatal(err)
	}
	onstatic.RegisterHandler(server.Mux)

	if conf.Variables.CGIMode {
		runCGIServerMode(ctx, server)
		return
	}

	runHTTPServerMode(ctx, server)
}

func runHTTPServerMode(ctx context.Context, server *http.Server) {
	log.Print("http server starting")
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
	if e := cgi.Serve(server.Mux); e != nil {
		log.Print(e)
	}

	server.Close()
}
