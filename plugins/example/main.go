package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	plugin "github.com/sters/onstatic/onstatic/plugin"
	"go.uber.org/zap"
)

type greeting struct{}

var _ plugin.API = (*greeting)(nil)

func (g *greeting) Initialize(context.Context) {}

func (g *greeting) Stop(context.Context) {}

func (g *greeting) Handlers() plugin.Handlers {
	return plugin.Handlers{
		"/greeting": func(res http.ResponseWriter, req *http.Request) {
			_, err := res.Write([]byte("Hello, greeting!"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%+v", err)
			}
		},
	}
}

// nolint
var EntryPoint = plugin.EntryPoint(func(context.Context, *zap.Logger) plugin.API {
	return &greeting{}
})
