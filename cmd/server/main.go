package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sters/staticman/conf"
	"github.com/sters/staticman/http"
	"github.com/sters/staticman/staticman"
)

func main() {
	conf.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := http.NewServer(conf.Variables.HTTPPort)
	if err != nil {
		panic(err)
	}
	staticman.RegisterHandler(server.Mux)

	log.Print("server starting")
	go func() { _ = server.Run() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	_ = server.Close()
}
