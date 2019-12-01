package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sters/staticman/http"
	"github.com/sters/staticman/staticman"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const port = ":18888"
	server, err := http.NewServer(port)
	if err != nil {
		panic(err)
	}
	staticman.RegisterHandler(server.Mux)

	go func() { _ = server.Run() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	_ = server.Close()
}
