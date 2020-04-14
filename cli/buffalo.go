package cli

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/clifix"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/bzr"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/fizz"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/flect"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/git"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/golang"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/grifts"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/i18n"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/mail"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/packr"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pkger"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/plush"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/refresh"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/soda"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/webpack"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugprint"
)

var _ plugcmd.SubCommander = &Buffalo{}
var _ plugins.Plugin = &Buffalo{}
var _ plugins.Scoper = &Buffalo{}
var _ plugprint.Describer = &Buffalo{}

// Buffalo represents the `buffalo` cli.
type Buffalo struct {
	plugins.Plugins
}

func NewFromRoot(root string) (*Buffalo, error) {
	b := &Buffalo{}

	b.Plugins = append(b.Plugins, clifix.Plugins()...)
	b.Plugins = append(b.Plugins, cmds.Plugins()...)
	b.Plugins = append(b.Plugins, fizz.Plugins()...)
	b.Plugins = append(b.Plugins, flect.Plugins()...)
	b.Plugins = append(b.Plugins, golang.Plugins()...)
	b.Plugins = append(b.Plugins, grifts.Plugins()...)
	b.Plugins = append(b.Plugins, i18n.Plugins()...)
	b.Plugins = append(b.Plugins, mail.Plugins()...)
	b.Plugins = append(b.Plugins, packr.Plugins()...)
	b.Plugins = append(b.Plugins, pkger.Plugins()...)
	b.Plugins = append(b.Plugins, plush.Plugins()...)
	b.Plugins = append(b.Plugins, pop.Plugins()...)
	b.Plugins = append(b.Plugins, refresh.Plugins()...)
	b.Plugins = append(b.Plugins, soda.Plugins()...)

	if _, err := os.Stat(filepath.Join(root, "package.json")); err == nil {
		b.Plugins = append(b.Plugins, webpack.Plugins()...)
	}

	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		b.Plugins = append(b.Plugins, git.Plugins()...)
	}

	if _, err := os.Stat(filepath.Join(root, ".bzr")); err == nil {
		b.Plugins = append(b.Plugins, bzr.Plugins()...)
	}

	sort.Slice(b.Plugins, func(i, j int) bool {
		return b.Plugins[i].PluginName() < b.Plugins[j].PluginName()
	})

	pfn := func() []plugins.Plugin {
		return b.Plugins
	}

	for _, b := range b.Plugins {
		f, ok := b.(plugins.Needer)
		if !ok {
			continue
		}
		f.WithPlugins(pfn)
	}

	return b, nil
}

func New() (*Buffalo, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return NewFromRoot(pwd)
}

func (b Buffalo) ScopedPlugins() []plugins.Plugin {
	return b.Plugins
}

func (b Buffalo) SubCommands() []plugins.Plugin {
	var plugs []plugins.Plugin
	for _, p := range b.ScopedPlugins() {
		if _, ok := p.(Commander); ok {
			plugs = append(plugs, p)
		}
	}
	return plugs
}

// Name ...
func (Buffalo) PluginName() string {
	return "buffalo"
}

func (Buffalo) String() string {
	return "buffalo"
}

// Description ...
func (Buffalo) Description() string {
	return "Tools for working with Buffalo applications"
}
