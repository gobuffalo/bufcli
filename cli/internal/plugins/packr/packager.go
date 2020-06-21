package packr

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/packr/v2/jam"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
)

var _ build.BeforeBuilder = &Packager{}
var _ build.Importer = &Packager{}
var _ build.AfterBuilder = &Packager{}
var _ build.Packager = &Packager{}

var _ plugcmd.Namer = &Packager{}
var _ plugins.Plugin = &Packager{}

var boxTemplate = `
package a
import "github.com/gobuffalo/packr/v2"

var (
	box = packr.New("buffalo:a:files", "./")
)

func findFile(name string) ([]byte, error) {
	return box.Find("database.yml")
}

`

type Packager struct{}

func (b *Packager) BeforeBuild(ctx context.Context, root string, args []string) error {
	err := os.RemoveAll(filepath.Join(root, "a"))
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(root, "a"), 0777)
	if err != nil {
		return err
	}

	return jam.Clean()
}

func (b *Packager) AfterBuild(ctx context.Context, root string, args []string, err error) error {
	return os.RemoveAll(filepath.Join(root, "a"))
}

func (b *Packager) Package(ctx context.Context, root string, files []string) error {
	err := b.copyFiles(root, files)
	if err != nil {
		return plugins.Wrap(b, err)
	}

	err = jam.Pack(jam.PackOptions{
		Roots: []string{root},
	})

	return plugins.Wrap(b, err)
}

func (b *Packager) BuildImports(ctx context.Context, root string) ([]string, error) {
	info, err := here.Current()
	if err != nil {
		return []string{}, plugins.Wrap(b, err)
	}

	return []string{filepath.Join(info.Module.Path, "a")}, nil
}

func (b *Packager) copyFiles(root string, files []string) error {

	for _, file := range files {
		f, err := ioutil.ReadFile(filepath.Join(file))
		if err != nil {
			return err
		}

		_, filename := filepath.Split(file)

		err = ioutil.WriteFile(filepath.Join(root, "a", filename), f, 0777)
		if err != nil {
			return err
		}
	}

	ioutil.WriteFile(filepath.Join(root, "a", "a.go"), []byte(boxTemplate), 0777)

	return nil
}

func (b Packager) PluginName() string {
	return "packr"
}

func (b Packager) CmdName() string {
	return "packr"
}
