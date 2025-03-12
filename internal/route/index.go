package route

import (
	"github.com/Penetration-Testing-Toolkit/ptt/internal/plugin"
	"github.com/Penetration-Testing-Toolkit/ptt/internal/templates"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
)

// IndexHandler "GET /" returns the whole root page.
func IndexHandler(c echo.Context, _ hclog.Logger, p []*plugin.ModulePlugin) error {
	i := plugin.PluginsToTemplate(p)

	return templates.Layout(i, templates.Index(i)).Render(c.Request().Context(), c.Response())
}
