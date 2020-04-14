package built

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/stretchr/testify/require"
)

type IniterFn func(ctx context.Context, root string, args []string) error

func (i IniterFn) BuiltInit(ctx context.Context, root string, args []string) error {
	return i(ctx, root, args)
}

func WithIniter(p plugins.Plugin, fn IniterFn) plugins.Plugin {
	type wi struct {
		IniterFn
		plugins.Plugin
	}
	return wi{
		Plugin:   p,
		IniterFn: fn,
	}
}

func Test_App_No_Args(t *testing.T) {
	r := require.New(t)

	var res bool
	app := &App{
		OriginalMain: func() {
			res = true
		},
	}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, "", args)
	r.NoError(err)
	r.True(res)
}

func Test_App_No_Args_EmptyApp(t *testing.T) {
	r := require.New(t)

	app := &App{}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, "", args)
	r.NoError(err)
}

func Test_App_Original_Args(t *testing.T) {
	r := require.New(t)

	var res bool
	app := &App{
		OriginalMain: func() {
			res = true
		},
	}

	args := []string{"--", "bob"}
	ctx := context.Background()
	err := app.Main(ctx, "", args)
	r.NoError(err)
	r.True(res)
}

func Test_App_No_Args_Fallthrough(t *testing.T) {
	r := require.New(t)

	var res bool
	app := &App{
		Fallthrough: func(ctx context.Context, root string, args []string) error {
			res = true
			return nil
		},
	}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, "", args)
	r.NoError(err)
	r.True(res)
}

func Test_App_With_Args_Fallthrough(t *testing.T) {
	r := require.New(t)

	var res bool
	app := &App{
		Fallthrough: func(ctx context.Context, root string, args []string) error {
			res = true
			return nil
		},
	}

	ctx := context.Background()
	err := app.Main(ctx, "", []string{"lee", "majors"})
	r.NoError(err)
	r.True(res)
}

func Test_App_With_Args_NoFallthrough(t *testing.T) {
	r := require.New(t)

	app := &App{}

	ctx := context.Background()
	root, err := os.Getwd()
	r.NoError(err)
	err = app.Main(ctx, root, []string{"lee", "majors"})
	r.NoError(err)
}

func Test_App_Init_Plugins(t *testing.T) {
	r := require.New(t)

	var res bool
	var pres bool

	fn := func(ctx context.Context, root string, args []string) error {
		pres = true
		return nil
	}

	plugs := plugins.Plugins{
		WithIniter(background(""), fn),
	}

	app := &App{
		OriginalMain: func() {
			res = true
		},
		Plugger: plugs,
	}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, "", args)
	r.NoError(err)
	r.True(res)
	r.True(pres)
}

func Test_App_Init_Plugins_Error(t *testing.T) {
	r := require.New(t)

	var res bool
	var pres bool
	exp := fmt.Errorf("boom")
	fn := func(ctx context.Context, root string, args []string) error {
		return exp
	}

	plugs := plugins.Plugins{
		WithIniter(background(""), fn),
	}

	app := &App{
		OriginalMain: func() {
			res = true
		},
		Plugger: plugs,
	}

	var args []string
	ctx := context.Background()
	err := app.Main(ctx, "", args)
	r.Error(err)
	r.Equal(exp, err)
	r.False(res)
	r.False(pres)
}

func Test_App_Version(t *testing.T) {
	r := require.New(t)

	stdout := &bytes.Buffer{}
	plugs := plugins.Plugins{
		plugio.NewOuter(stdout),
	}
	app := &App{
		BuildVersion: "v1",
		Plugger:      plugs,
	}
	ctx := context.Background()

	err := app.Main(ctx, "", []string{"version"})
	r.NoError(err)

	s := strings.TrimSpace(stdout.String())
	r.Equal("v1", s)
}
