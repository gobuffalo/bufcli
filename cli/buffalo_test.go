package cli

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertPluginContained(name string, plugins plugins.Plugins) assert.Comparison {
	return func() (success bool) {
		success = false
		for _, p := range plugins {
			if p.PluginName() == name {
				return true
			}
		}
		return
	}
}

func Test_Buffalo_New(t *testing.T) {
	r := require.New(t)

	b, err := New()
	r.NoError(err)
	r.NotNil(b)
	r.NotEmpty(b.Plugins)
}

func Test_Buffalo_New_WithFiles(t *testing.T) {
	r := require.New(t)

	root, err := os.Getwd()
	r.NoError(err)

	//create files
	r.NoError(ioutil.WriteFile(filepath.Join(root, "package.json"), []byte(""), 0600))
	r.NoError(ioutil.WriteFile(filepath.Join(root, ".git"), []byte(""), 0600))
	r.NoError(ioutil.WriteFile(filepath.Join(root, ".bzr"), []byte(""), 0600))

	b, err := New()
	r.NoError(err)
	r.NotNil(b)
	r.NotEmpty(b.Plugins)

	//test for webpack builder include
	r.Condition(AssertPluginContained("webpack/builder", b.Plugins))

	r.NoError(os.Remove("package.json"))

	//test for git builder include
	r.Condition(AssertPluginContained("git/versioner", b.Plugins))

	r.NoError(os.Remove(".git"))

	//test for bzr builder include
	r.Condition(AssertPluginContained("bzr", b.Plugins))

	r.NoError(os.Remove(".bzr"))
}

func Test_Buffalo_SubCommands(t *testing.T) {
	r := require.New(t)

	c := &cp{}
	b := &Buffalo{
		Plugins: plugins.Plugins{
			background("foo"),
			c,
		},
	}
	r.Len(b.Plugins, 2)

	cmds := b.SubCommands()
	r.Len(cmds, 1)
	r.Equal(c, cmds[0])
}
