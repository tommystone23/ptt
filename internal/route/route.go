package route

import (
	"github.com/Penetration-Testing-Toolkit/ptt/internal/plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

type HandlerFunc func(echo.Context, hclog.Logger, []*plugin.ModulePlugin) error
