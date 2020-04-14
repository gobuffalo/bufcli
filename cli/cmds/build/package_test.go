package build

import (
	"context"
	"errors"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_Cmd_Package(t *testing.T) {
	r := require.New(t)

	pkg := &packager{
		files: []string{"A"},
	}
	pf := &packFiler{
		files: []string{"B"},
	}

	plugs := plugins.Plugins{
		pkg,
		pf,
		&bladeRunner{},
	}

	bc := &Cmd{}
	bc.WithPlugins(func() []plugins.Plugin {
		return plugs
	})

	err := bc.Main(context.Background(), ".", nil)
	r.NoError(err)

	r.Len(pkg.files, 2)
	r.Equal([]string{"A", "B"}, pkg.files)
}

func Test_Cmd_Package_Err(t *testing.T) {
	r := require.New(t)

	//Test error from PackFiler
	pkg := &packager{
		files: []string{"A"},
	}
	pf := &packFiler{
		files: []string{"B"},
		err:   errors.New("Bad PackFiler"),
	}

	plugs := plugins.Plugins{
		pkg,
		pf,
		&bladeRunner{},
	}

	bc := &Cmd{}
	bc.WithPlugins(func() []plugins.Plugin {
		return plugs
	})

	err := bc.Main(context.Background(), ".", nil)
	r.Error(err)
	r.Contains(err.Error(), "Bad PackFiler")

	//Test error from Packager
	pkg = &packager{
		files: []string{"A"},
		err:   errors.New("Bad Packager"),
	}
	pf = &packFiler{
		files: []string{"B"},
	}

	plugs = plugins.Plugins{
		pkg,
		pf,
		&bladeRunner{},
	}

	bc = &Cmd{}
	bc.WithPlugins(func() []plugins.Plugin {
		return plugs
	})

	err = bc.Main(context.Background(), ".", nil)
	r.Error(err)
	r.Contains(err.Error(), "Bad Packager")
}
