package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/stretchr/testify/require"
)

func Test_Buffalo_Help(t *testing.T) {
	r := require.New(t)

	stdout := &bytes.Buffer{}

	b := &Buffalo{
		Plugins: plugins.Plugins{
			plugio.NewOuter(stdout),
			TestNeeder{
				Testy:               t,
				ExpectedPluginCount: 5,
			},
			TestOutNeeder{},
			TestInNeeder{},
			TestErrNeeder{},
		},
	}

	ctx := context.Background()

	args := []string{"-h"}

	err := b.Main(ctx, "", args)
	r.NoError(err)

	r.Contains(stdout.String(), b.Description())
}

func Test_Buffalo_FailingPlugins(t *testing.T) {
	r := require.New(t)

	stdout := &bytes.Buffer{}

	//Bad Out
	b := &Buffalo{
		Plugins: plugins.Plugins{
			plugio.NewOuter(stdout),
			TestNeeder{
				Testy:               t,
				ExpectedPluginCount: 5,
			},
			TestOutNeeder{Error: errors.New("Bad Out")},
			TestInNeeder{},
			TestErrNeeder{},
		},
	}

	ctx := context.Background()

	args := []string{"-h"}

	err := b.Main(ctx, "", args)
	r.Error(err)
	r.Equal("Bad Out", err.Error())

	//Bad In
	b = &Buffalo{
		Plugins: plugins.Plugins{
			TestInNeeder{Error: errors.New("Bad In")},
		},
	}

	err = b.Main(ctx, "", args)
	r.Error(err)
	r.Equal("Bad In", err.Error())

	//Bad Err
	b = &Buffalo{
		Plugins: plugins.Plugins{
			TestErrNeeder{Error: errors.New("Bad Err")},
		},
	}

	err = b.Main(ctx, "", args)
	r.Error(err)
	r.Equal("Bad Err", err.Error())

	//Bad Combo
	b = &Buffalo{
		Plugins: plugins.Plugins{
			TestNeeder{
				Testy:               t,
				ExpectedPluginCount: 1,
				Error:               errors.New("Bad Combo"),
			},
		},
	}

	err = b.Main(ctx, "", args)
	r.Error(err)
	r.Equal("Bad Combo", err.Error())
}

func Test_Buffalo_Main_SubCommand(t *testing.T) {
	r := require.New(t)

	c := &cp{}
	b := &Buffalo{
		Plugins: plugins.Plugins{
			c,
		},
	}

	ctx := context.Background()

	args := []string{c.PluginName()}

	exp := []string{"hello"}
	args = append(args, exp...)

	err := b.Main(ctx, "", args)
	r.NoError(err)
	r.Equal(exp, c.args)
}

func Test_Buffalo_Main_SubCommand_Alias(t *testing.T) {
	r := require.New(t)

	c := &cp{aliases: []string{"sc"}}
	b := &Buffalo{
		Plugins: plugins.Plugins{
			c,
		},
	}

	ctx := context.Background()

	args := []string{"sc"}

	exp := []string{"hello"}
	args = append(args, exp...)

	err := b.Main(ctx, "", args)
	r.NoError(err)
	r.Equal(exp, c.args)
}

func Test_Buffalo_WithoutPlugins(t *testing.T) {
	r := require.New(t)

	b := &Buffalo{}

	ctx := context.Background()

	args := []string{"sc"}

	exp := []string{"hello"}
	args = append(args, exp...)

	err := b.Main(ctx, "", args)
	r.NoError(err)
}
