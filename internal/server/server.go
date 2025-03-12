package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Penetration-Testing-Toolkit/ptt/internal/plugin"
	"github.com/Penetration-Testing-Toolkit/ptt/internal/route"
	"github.com/Penetration-Testing-Toolkit/ptt/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	Logger  hclog.Logger
	Static  fs.FS
	Address string
}

func Start(cfg *Config) {
	logger := cfg.Logger

	// Create Echo server
	e := echo.New()

	// Format echo's logging through hashicorp/go-hclog
	requestLogger := logger.Named("request")
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP: true,
		LogURI:      true,
		LogMethod:   true,
		LogProtocol: true,
		LogStatus:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			requestLogger.With("IP", v.RemoteIP).
				With("URI", v.URI).
				With("method", v.Method).
				With("protocol", v.Protocol).
				With("status", v.Status).
				Info("new request")
			return nil
		},
	}))

	// Custom error handler
	errHandler := func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
		}
		logger.Error(err.Error())
		err = c.String(code, "something went wrong, internal server error")
		if err != nil {
			return
		}
		return
	}
	e.HTTPErrorHandler = errHandler

	// Host static assets
	staticServer := http.FileServer(http.FS(cfg.Static))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticServer)))
	e.GET("/favicon.ico", echo.WrapHandler(staticServer))

	// Start plugin discovery & register routes
	shared.Logger = logger
	plugins, cleanup := plugin.StartPlugins(logger, e)
	defer cleanup(plugins)

	// Setup server's routes
	setupRoutes(logger, e, plugins)

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if logger.GetLevel() <= hclog.Debug {
			for _, r := range e.Routes() {
				logger.Debug("route defined", "method", fmt.Sprintf("%v", r.Method),
					"path", fmt.Sprintf("%v", r.Path))
			}
		}

		if err := e.Start(cfg.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("shutting down the server")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Error("error shutting down server", "err", err.Error())
		os.Exit(1)
	}
}

func setupRoutes(l hclog.Logger, e *echo.Echo, p []*plugin.ModulePlugin) {
	// Adapts our handler style (route.HandlerFunc) into an echo.HandlerFunc
	// This allows us to pass the hclog.Logger and plugins alongside echo.Context
	adapter := func(l hclog.Logger, p []*plugin.ModulePlugin, f route.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return f(c, l, p)
		}
	}

	e.GET("/", adapter(l, p, route.IndexHandler))
}
