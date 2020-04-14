package build

import (
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Subcommands(t *testing.T) {
	r := require.New(t)

	b := &builder{}
	all := plugins.Plugins{
		background("foo"),
		&beforeBuilder{},
		b,
		&afterBuilder{},
		background("bar"),
		&buildVersioner{},
		&packager{},
		&bladeRunner{},
	}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return all
		},
	}

	plugs := bc.SubCommands()
	r.Len(plugs, 1)
	r.Equal(b, plugs[0])
}

func Test_Cmd_ScopedPlugins(t *testing.T) {
	r := require.New(t)

	all := plugins.Plugins{
		background("foo"),
		&builder{},
		&beforeBuilder{},
		&afterBuilder{},
		background("bar"),
		&buildVersioner{},
		&buildImporter{},
		&bladeRunner{},
		&packager{},
		&buildTagger{},
		&buildStdouter{},
	}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return all
		},
	}

	plugs := bc.ScopedPlugins()
	r.NotEqual(all, plugs)

	ep := plugins.Plugins(plugs)

	tot := len(all) - 2
	r.Equal(tot, len(ep))

}
