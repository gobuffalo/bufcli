package build

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_GoCmd(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	cli := filepath.Join("bin", "build")
	if runtime.GOOS == "windows" {
		cli += ".exe"
	}
	exp := []string{"go", "build", "-o", cli}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_Bin(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin: "cli",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_Mod(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin: "cli",
		mod: "vendor",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-mod", "vendor"}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_Tags(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin:  "cli",
		tags: "a b c",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-tags", "a b c"}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_Static(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin:    "cli",
		static: true,
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-ldflags", "-linkmode external -extldflags \"-static\""}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_LDFlags(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{
		bin:     "cli",
		ldFlags: "linky",
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)

	n := "cli"
	if runtime.GOOS == "windows" {
		n = "cli.exe"
	}

	exp := []string{"go", "build", "-o", n, "-ldflags", "linky"}
	r.Equal(exp, cmd.Args)
}

func Test_Cmd_GoCmd_BadRoot(t *testing.T) {
	r := require.New(t)

	bc := &Cmd{}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, "")
	r.Error(err)
	r.Nil(cmd)
}

func Test_Cmd_GoCmd_Tagger(t *testing.T) {
	r := require.New(t)

	pfn := func() []plugins.Plugin {
		return []plugins.Plugin{
			&buildTagger{},
		}
	}

	bc := &Cmd{
		pluginsFn: pfn,
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.NoError(err)
	r.NotNil(cmd)
}

func Test_Cmd_GoCmd_TaggerErr(t *testing.T) {
	r := require.New(t)

	pfn := func() []plugins.Plugin {
		return []plugins.Plugin{
			&buildTagger{err: errors.New("Bad Tagger")},
		}
	}

	bc := &Cmd{
		pluginsFn: pfn,
	}

	ctx := context.Background()
	cmd, err := bc.GoCmd(ctx, ".")
	r.Error(err)
	r.Nil(cmd)
	r.Contains(err.Error(), "Bad Tagger")
}

func Test_Cmd_GoBuild_Build_Err(t *testing.T) {
	r := require.New(t)

	pfn := func() []plugins.Plugin {
		return []plugins.Plugin{
			&bladeRunner{err: errors.New("Bad Runner")},
		}
	}

	bc := &Cmd{
		pluginsFn: pfn,
	}

	ctx := context.Background()
	var args []string
	err := bc.build(ctx, "", args)
	r.Error(err)
	if runtime.GOOS == "windows" {
		r.Contains(err.Error(), "The system cannot find the path specified")
	} else {
		r.Contains(err.Error(), "no such file or directory")
	}

	err = bc.build(ctx, ".", args)
	r.Error(err)
	r.Contains(err.Error(), "Bad Runner")
}
