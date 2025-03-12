package main

import (
	"embed"
	"github.com/Penetration-Testing-Toolkit/ptt/internal/server"
	"github.com/Penetration-Testing-Toolkit/ptt/shared"
	"github.com/hashicorp/go-hclog"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

//go:embed static
var embeddedFS embed.FS

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
		Name:       "app",
		Level:      logLevel,
		Output:     os.Stdout,
		JSONFormat: json,
	})
	logger.Debug("logger created", "level", logLevel, "json", json)

	// Determine if files will be hosted from embedded in the binary or from the OS file system
	var f fs.FS
	env := os.Getenv("ENV")
	if strings.ToUpper(env) == "DEV" {
		// If dev environment, host static assets from file system
		f = os.DirFS("static")
		logger.Debug("hosting static files from file system")
	} else {
		// Else, host static assets embedded in Go binary
		f, err = fs.Sub(embeddedFS, "static")
		if err != nil {
			logger.Error("error finding 'static' directory in embedded filesystem", "err", err.Error())
			panic(err)
		}
		logger.Debug("hosting static files from embedded file system")
	}

	server.Start(&server.Config{
		Logger:  logger,
		Static:  f,
		Address: ":8080",
	})
}
