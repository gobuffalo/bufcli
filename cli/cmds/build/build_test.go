package build

import (
	"path/filepath"
	"testing"

	"github.com/gobuffalo/here"
	"github.com/stretchr/testify/require"
)

type background string

func (b background) PluginName() string {
	return string(b)
}

func newRef(t *testing.T, root string) here.Info {
	t.Helper()

	info := here.Info{
		Dir:        root,
		ImportPath: "github.com/markbates/coke",
		Name:       "coke",
		Root:       root,
		Module: here.Module{
			Path:  "github.com/markbates/coke",
			Main:  true,
			Dir:   root,
			GoMod: filepath.Join(root, "go.mod"),
		},
	}

	return info
}

func Test_Plugins(t *testing.T) {
	r := require.New(t)
	plugs := Plugins()
	r.Len(plugs, 2)
	for _, p := range plugs {
		if p.PluginName() != "main" && p.PluginName() != "build" {
			r.FailNow("should only be main and build")
		}
	}
}
