package build

import (
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_FindBuilder(t *testing.T) {
	r := require.New(t)

	bb := &buildBuilder{}
	plugs := []plugins.Plugin{
		&buildFlagger{},
		bb,
	}

	builder := FindBuilder("buildBuilder", plugs)
	r.NotNil(builder)
	r.Equal(bb.PluginName(), builder.PluginName())
}

func Test_FindBuilder_NoBuilder(t *testing.T) {
	r := require.New(t)

	plugs := []plugins.Plugin{
		&buildFlagger{},
		&buildImporter{},
	}

	builder := FindBuilder("buildBuilder", plugs)
	r.Nil(builder)
}

func Test_FindBuilder_With_Namer(t *testing.T) {
	r := require.New(t)

	bn := &buildNamer{cmdName: "first"}
	plugs := []plugins.Plugin{
		&buildFlagger{},
		bn,
	}

	expects := []struct {
		Term    string
		IsNil   bool
		IsNamer bool
	}{
		{"first", false, true},
		{"second", true, false},
		{"buildNamer", false, true},
		{"third", true, false},
		{"buildBuilder", true, false},
	}

	for _, e := range expects {
		builder := FindBuilder(e.Term, plugs)
		if e.IsNil {
			r.Nil(builder)
		} else {
			r.NotNil(builder)
			r.Equal(bn.PluginName(), builder.PluginName())
			namer, ok := builder.(Namer)
			if e.IsNamer {
				r.True(ok)
				r.Equal(bn.CmdName(), namer.CmdName())
			} else {
				r.False(ok)
			}
		}
	}
}

func Test_FindBuilder_With_Aliaser(t *testing.T) {
	r := require.New(t)

	ba := &buildAliaser{aliases: []string{"first", "second"}}
	plugs := []plugins.Plugin{
		&buildFlagger{},
		ba,
	}

	expects := []struct {
		Term      string
		IsNil     bool
		IsAliaser bool
	}{
		{"first", false, true},
		{"second", false, true},
		{"buildAliaser", false, true},
		{"third", true, false},
	}

	for _, expect := range expects {
		builder := FindBuilder(expect.Term, plugs)
		if expect.IsNil {
			r.Nil(builder)
		} else {
			r.NotNil(builder)
			r.Equal(ba.PluginName(), builder.PluginName())
			aliaser, ok := builder.(Aliaser)
			if expect.IsAliaser {
				r.True(ok)
				r.Equal(ba.CmdAliases(), aliaser.CmdAliases())
			} else {
				r.False(ok)
			}
		}
	}
}
