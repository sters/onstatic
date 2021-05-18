package main

import (
	"fmt"
	"net/http"
	"os"

	plugin "github.com/sters/onstatic/onstatic/plugin"
	"go.uber.org/zap"
)

type greeting struct{}

var _ plugin.API = (*greeting)(nil)

func (g *greeting) Register() plugin.Handlers {
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
var EntryPoint = plugin.EntryPoint(func(*zap.Logger) plugin.API {
	return &greeting{}
})
