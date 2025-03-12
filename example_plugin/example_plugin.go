package main

import (
	"context"
	"fmt"
	"github.com/Penetration-Testing-Toolkit/ptt/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var info = &shared.ModuleInfo{
	ID:      "github.com/chronotrax/example_plugin",
	Name:    "Example Plugin",
	Version: "1.1.0",
}

type HandlerFunc func(context.Context, *http.Request) (*shared.Response, error)
type SSEHandlerFunc func(context.Context, *http.Request) (chan *shared.Response, error)

// Router is a simple router to match a request's method & path to the correct handler function.
type Router struct {
	// Method -> Path -> HandlerFunc
	routes map[string]map[string]HandlerFunc

	// Method -> Path -> SSEHandlerFunc
	sseRoutes map[string]map[string]SSEHandlerFunc
}

var router = &Router{
	routes:    make(map[string]map[string]HandlerFunc),
	sseRoutes: make(map[string]map[string]SSEHandlerFunc),
}

// ModuleExample is a real implementation of a shared.Module plugin.
// It uses hclog.Logger for logging to the hashicorp/go-plugin system.
type ModuleExample struct {
	logger hclog.Logger
}

// Register implements shared.Module's Register.
func (m *ModuleExample) Register(_ context.Context) (*shared.ModuleInfo, error) {
	return info, nil
}

// Handle implements shared.Module's Handle.
func (m *ModuleExample) Handle(ctx context.Context, req *http.Request) (*shared.Response, error) {
	// Lookup HandlerFunc in router
	handler, exists := router.routes[req.Method][req.URL.String()]
	if !exists {
		return nil, fmt.Errorf("handler does not exist for route %v %v", req.Method, req.URL)
	}

	return handler(ctx, req)
}

// HandleSSE implements shared.Module's HandleSSE.
func (m *ModuleExample) HandleSSE(ctx context.Context, req *http.Request) (chan *shared.Response, error) {
	// Lookup SSEHandlerFunc in router
	sseHandler, exists := router.sseRoutes[req.Method][req.URL.String()]
	if !exists {
		return nil, fmt.Errorf("sse sseHandler does not exist for route %v %v", req.Method, req.URL)
	}

	return sseHandler(ctx, req)
}

func main() {
	// Read environment variables to configure logging
	json := true
	if strings.ToUpper(os.Getenv("JSON")) == "FALSE" {
		json = false
	}
	logLevel := shared.LoggerOptions.Level
	l := os.Getenv("LOG_LEVEL")
	i, err := strconv.Atoi(l)
	if err == nil {
		logLevel = hclog.Level(i)
	}

	// Create a hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:        info.ID,
		Level:       logLevel,
		Output:      shared.LoggerOptions.Output,
		JSONFormat:  json,
		DisableTime: shared.LoggerOptions.DisableTime,
	})

	// Create the ModuleExample instance
	module := &ModuleExample{
		logger: logger,
	}

	// Setup plugin's routes
	module.setupRoutes()

	// Setup hashicorp/go-plugin stuff
	shared.Logger = logger
	shared.PluginMap["module"] = &shared.ModuleGRPCPlugin{Impl: module}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins:         shared.PluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
		Logger:          logger,
	})
}
