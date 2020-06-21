package cmds

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/develop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/fix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/generate"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/info"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/newapp"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/resource"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/setup"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/test"
	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/version"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/packr"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/refresh"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/soda"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack"
	"github.com/gobuffalo/buffalo-cli/v2/meta"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, insidePlugins()...)
	plugs = append(plugs, outsidePlugins()...)
	plugs = append(plugs, version.Plugins()...)
	return plugs
}

func AvailablePlugins(root string) []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, version.Plugins()...)
	if meta.IsBuffalo(root) {
		plugs = append(plugs, insidePlugins()...)
		return plugs
	}
	plugs = append(plugs, outsidePlugins()...)
	return plugs
}

func outsidePlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, info.Plugins()...)
	plugs = append(plugs, newapp.Plugins()...)
	return plugs
}

func insidePlugins() []plugins.Plugin {
	var plugs []plugins.Plugin
	plugs = append(plugs, build.Plugins()...)
	plugs = append(plugs, develop.Plugins()...)
	plugs = append(plugs, fix.Plugins()...)
	plugs = append(plugs, generate.Plugins()...)
	plugs = append(plugs, info.Plugins()...)
	plugs = append(plugs, resource.Plugins()...)
	plugs = append(plugs, setup.Plugins()...)
	plugs = append(plugs, test.Plugins()...)

	plugs = append(plugs, webpack.Plugins()...)
	plugs = append(plugs, refresh.Plugins()...)
	plugs = append(plugs, soda.Plugins()...)
	plugs = append(plugs, packr.Plugins()...)
	plugs = append(plugs, pop.Plugins()...)

	return plugs
}
