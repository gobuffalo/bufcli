package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Main(t *testing.T) {
	r := require.New(t)

	info := newRef(t, "")

	bc := &Cmd{}
	bc.WithHereInfo(info)

	bn := filepath.Join("bin", "coke")
	if runtime.GOOS == "windows" {
		bn += ".exe"
	}
	exp := []string{"go", "build", "-o", bn}

	br := &bladeRunner{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{br}
	})

	var args []string
	err := bc.Main(context.Background(), args)
	r.NoError(err)
	r.NotNil(br.cmd)
	r.Equal(exp, br.cmd.Args)
}

func Test_Cmd_Main_SubCommand(t *testing.T) {
	r := require.New(t)

	p := &builder{name: "foo"}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{p.name, "a", "b", "c"}

	err := bc.Main(context.Background(), args)
	r.NoError(err)
	r.Equal([]string{"a", "b", "c"}, p.args)
}

func Test_Cmd_Main_SubCommand_err(t *testing.T) {
	r := require.New(t)

	p := &builder{name: "foo", err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{p.name, "a", "b", "c"}

	err := bc.Main(context.Background(), args)
	r.Error(err)
}

func Test_Cmd_Main_ValidateTemplates(t *testing.T) {
	r := require.New(t)

	p := &templatesValidator{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	info := newRef(t, "")

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}
	bc.WithHereInfo(info)

	args := []string{}

	err := bc.Main(context.Background(), args)
	r.NoError(err)

	bi, err := bc.HereInfo()
	r.NoError(err)
	r.Equal(bi.Root, p.root)
}

func Test_Cmd_Main_ValidateTemplates_err(t *testing.T) {
	r := require.New(t)

	p := &templatesValidator{err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{}

	err := bc.Main(context.Background(), args)
	r.Error(err)
}

func Test_Cmd_Main_BeforeBuilders(t *testing.T) {
	r := require.New(t)

	p := &beforeBuilder{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), args)
	r.NoError(err)
}

func Test_Cmd_Main_BeforeBuilders_err(t *testing.T) {
	r := require.New(t)

	p := &beforeBuilder{err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), args)
	r.Error(err)
}

func Test_Cmd_Main_AfterBuilders(t *testing.T) {
	r := require.New(t)

	p := &afterBuilder{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), args)
	r.NoError(err)
}

func Test_Cmd_Main_AfterBuilders_err(t *testing.T) {
	r := require.New(t)

	b := &beforeBuilder{err: fmt.Errorf("error")}
	a := &afterBuilder{}
	plugs := plugins.Plugins{a, b, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	var args []string

	err := bc.Main(context.Background(), args)
	r.Error(err)
	r.Equal(err, a.err)
}
