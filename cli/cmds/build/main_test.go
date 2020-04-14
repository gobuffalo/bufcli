package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Main(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	bn := filepath.Join("bin", "build")
	if runtime.GOOS == "windows" {
		bn += ".exe"
	}
	exp := []string{"go", "build", "-o", bn}

	br := &bladeRunner{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{br}
	})

	var args []string
	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
	r.NotNil(br.cmd)
	r.Equal(exp, br.cmd.Args)

	args = []string{"--help"}
	err = bc.Main(context.Background(), ".", args)
	r.NoError(err)
	r.NotNil(br.cmd)
	r.Equal(exp, br.cmd.Args)
}

func Test_Cmd_Main_FlagErr(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	bn := filepath.Join("bin", "build")
	if runtime.GOOS == "windows" {
		bn += ".exe"
	}

	br := &bladeRunner{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{br}
	})

	args := []string{"-help"}
	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
}

func Test_Cmd_Main_BadRootErr(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	bn := filepath.Join("bin", "build")
	if runtime.GOOS == "windows" {
		bn += ".exe"
	}

	br := &bladeRunner{}
	bc.WithPlugins(func() []plugins.Plugin {
		return []plugins.Plugin{br}
	})

	args := []string{}
	err := bc.Main(context.Background(), "", args)
	r.Error(err)
	if runtime.GOOS == "windows" {
		r.Contains(err.Error(), "The system cannot find the path specified")
	} else {
		r.Contains(err.Error(), "no such file or directory")
	}
}

func Test_Cmd_Main_SubCommand(t *testing.T) {
	r := require.New(t)

	p := &builder{name: "foo"}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: plugs.ScopedPlugins,
	}

	args := []string{p.name, "a", "b", "c"}

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
	r.Equal([]string{"a", "b", "c"}, p.args)
}

func Test_Cmd_Main_SubCommand_err(t *testing.T) {
	r := require.New(t)

	p := &builder{name: "foo", err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	args := []string{p.name, "a", "b", "c"}

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
}

func Test_Cmd_Main_BeforeBuilders(t *testing.T) {
	r := require.New(t)

	p := &beforeBuilder{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
}

func Test_Cmd_Main_BeforeBuilders_err(t *testing.T) {
	r := require.New(t)

	p := &beforeBuilder{err: fmt.Errorf("error")}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
}

func Test_Cmd_Main_AfterBuilders(t *testing.T) {
	r := require.New(t)

	p := &afterBuilder{}
	plugs := plugins.Plugins{p, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
}

func Test_Cmd_Main_AfterBuilders_err(t *testing.T) {
	r := require.New(t)

	b := &beforeBuilder{err: fmt.Errorf("science fiction twin")}
	a := &afterBuilder{}
	plugs := plugins.Plugins{a, b, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	var args []string

	err := bc.Main(context.Background(), ".", args)
	r.Error(err)
	r.Contains(err.Error(), b.err.Error())

	b = &beforeBuilder{}
	a = &afterBuilder{err: fmt.Errorf("science fiction twin")}
	plugs = plugins.Plugins{a, b, &bladeRunner{}}

	bc = &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	err = bc.Main(context.Background(), ".", args)
	r.Error(err)
	r.Contains(err.Error(), a.err.Error())
}

func Test_Cmd_Main_No_AfterBuilders_Err_On_Help(t *testing.T) {
	r := require.New(t)

	b := &beforeBuilder{}
	a := &afterBuilder{err: fmt.Errorf("science fiction twin")}
	plugs := plugins.Plugins{a, b, &bladeRunner{}}

	bc := &Cmd{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	args := []string{"--help"}

	err := bc.Main(context.Background(), ".", args)
	r.NoError(err)
}
