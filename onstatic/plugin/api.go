package plugin

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	PluginDir                = ".onstatic"
	PluginExportVariableName = "EntryPoint"
)

type Handler func(res http.ResponseWriter, req *http.Request)

type Endpoint string

type Handlers map[Endpoint]Handler

type API interface {
	Register() Handlers
}

type EntryPoint func(*zap.Logger) API
