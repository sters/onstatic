package plugin

import (
	"net/http"
)

const (
	PluginFilePath           = ".onstatic/plugin.so"
	PluginExportVariableName = "API"
)

type Handler func(res http.ResponseWriter, req *http.Request)

type Endpoint string

type Handlers map[Endpoint]Handler

type API interface {
	Register() Handlers
}
