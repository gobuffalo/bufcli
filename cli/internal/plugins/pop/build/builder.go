package build

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/cmds/build"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/refresh"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/pop/v5/soda/cmd"
)

var _ build.Importer = Builder{}
var _ build.PackFiler = &Builder{}
var _ refresh.Tagger = &Builder{}
var _ build.BuildArger = &Builder{}
var _ build.Versioner = &Builder{}
var _ plugins.Plugin = Builder{}

const filePath = "/database.yml"

var initCode = `
package a

import (
	"bytes"
	"log"

	"github.com/gobuffalo/pop/v5"
)

func init() {
	//findFile should be provided by packagers
	file, err := findFile("database.yml")
	if err != nil {
		log.Fatal(err)
	}
	  
	r := bytes.NewReader(file)
	err = pop.LoadFrom(r)
	if err != nil {
		log.Fatal(err)
	}
}
`

type Builder struct{}

func (Builder) PluginName() string {
	return "pop/builder"
}

func (bd *Builder) GoBuildArgs(ctx context.Context, root string, args []string) ([]string, error) {
	x := []string{}

	tags, err := bd.RefreshTags(ctx, root)
	if err != nil || len(tags) == 0 {
		return x, plugins.Wrap(bd, err)
	}

	if len(tags) == 0 {
		return x, nil
	}

	x = append(x, []string{"-tags", strings.Join(tags, " ")}...)
	return x, nil
}

func (bd *Builder) BeforeBuild(ctx context.Context, root string, args []string) error {
	return ioutil.WriteFile(filepath.Join(root, "a", "pop.go"), []byte(initCode), 0777)
}

func (bd *Builder) RefreshTags(ctx context.Context, root string) ([]string, error) {
	var args []string
	dy := filepath.Join(root, "database.yml")
	if _, err := os.Stat(dy); err != nil {
		return args, nil
	}

	b, err := ioutil.ReadFile(dy)
	if err != nil {
		return nil, plugins.Wrap(bd, err)
	}
	if bytes.Contains(b, []byte("sqlite")) {
		args = append(args, "sqlite")
	}
	return args, nil
}

func (bd *Builder) BuildVersion(ctx context.Context, root string) (string, error) {
	return cmd.Version, nil
}

func (bd *Builder) PackageFiles(ctx context.Context, root string) ([]string, error) {
	return []string{
		filepath.Join(root, filePath),
	}, nil
}

func (Builder) BuildImports(ctx context.Context, root string) ([]string, error) {
	dir := filepath.Join(root, "models")
	info, err := here.Dir(dir)
	if err != nil {
		return nil, nil
	}
	return []string{info.ImportPath}, nil
}
