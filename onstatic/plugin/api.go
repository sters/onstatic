package plugin

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

const (
	PluginDir                = ".onstatic"
	PluginExportVariableName = "EntryPoint"
)

// Endpoint definition that should start "/"
type Endpoint string

// Handlers is Endpoint-HandlerFunc collection
type Handlers map[Endpoint]http.HandlerFunc

// API is main structure of this plugin
type API interface {
	// Initialize this API
	Initialize(context.Context)
	// Start this API handling
	Stop(context.Context)
	// Handlers returns it for this API
	Handlers() Handlers
}

// EntryPoint is plugin entry point. First, call this function.
type EntryPoint func(context.Context, *zap.Logger) API
