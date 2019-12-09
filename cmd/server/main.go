package main

import (
	"context"
	"log"
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
		panic(err)
	}
	onstatic.RegisterHandler(server.Mux)

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
