package build

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

var _ plugcmd.Aliaser = &Cmd{}
var _ plugcmd.SubCommander = &Cmd{}
var _ plugins.Needer = &Cmd{}
var _ plugins.Plugin = &Cmd{}
var _ plugins.Scoper = &Cmd{}
var _ plugprint.Describer = &Cmd{}
var _ plugprint.FlagPrinter = &Cmd{}

type Cmd struct {
	// Mod is the -mod flag
	mod string
	// Static sets the following flags for the final `go build` command:
	// -linkmode external
	// -extldflags "-static"
	static bool
	// Environment the binary is meant for. defaults to "development"
	environment string
	// LDFlags to be passed to the final `go build` command
	ldFlags string
	// BuildFlags to be passed to the final `go build` command
	buildFlags []string
	tags       string
	bin        string
	verbose    bool

	help      bool
	pluginsFn plugins.Feeder
	flags     *pflag.FlagSet
}

func (cmd *Cmd) WithPlugins(f plugins.Feeder) {
	cmd.pluginsFn = f
}

func (*Cmd) CmdAliases() []string {
	return []string{"b", "install"}
}

func (cmd Cmd) PluginName() string {
	return "build"
}

func (Cmd) Description() string {
	return "Build the application binary, including bundling of webpack (packr & webpack)"
}

func (cmd *Cmd) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range cmd.ScopedPlugins() {
		if _, ok := p.(Builder); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

func (cmd *Cmd) ScopedPlugins() []plugins.Plugin {
	if cmd.pluginsFn == nil {
		return nil
	}

	var plugs []plugins.Plugin
	for _, p := range cmd.pluginsFn() {
		switch p.(type) {
		case BuildArger:
			plugs = append(plugs, p)
		case Builder:
			plugs = append(plugs, p)
		case BeforeBuilder:
			plugs = append(plugs, p)
		case AfterBuilder:
			plugs = append(plugs, p)
		case Versioner:
			plugs = append(plugs, p)
		case Packager:
			plugs = append(plugs, p)
		case PackFiler:
			plugs = append(plugs, p)
		case Flagger:
			plugs = append(plugs, p)
		case Pflagger:
			plugs = append(plugs, p)
		case Importer:
			plugs = append(plugs, p)
		case GoBuilder:
			plugs = append(plugs, p)
		case Stdouter:
			plugs = append(plugs, p)
		case Stdiner:
			plugs = append(plugs, p)
		case Stderrer:
			plugs = append(plugs, p)
		}
	}
	return plugs
}
