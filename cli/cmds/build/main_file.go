package build

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugprint"
	"golang.org/x/tools/go/ast/astutil"
)

const mainBuildFile = "main.build.go"

var _ AfterBuilder = &MainFile{}
var _ BeforeBuilder = &MainFile{}
var _ plugins.Plugin = &MainFile{}
var _ plugins.Needer = &MainFile{}
var _ plugins.Scoper = &MainFile{}
var _ plugprint.Hider = &MainFile{}

type MainFile struct {
	pluginsFn         plugins.Feeder
	withFallthroughFn func() bool
}

func (bc *MainFile) WithPlugins(f plugins.Feeder) {
	bc.pluginsFn = f
}

func (MainFile) HidePlugin() {}

func (MainFile) PluginName() string {
	return "main"
}

func (bc *MainFile) ScopedPlugins() []plugins.Plugin {
	if bc.pluginsFn == nil {
		return nil
	}
	return bc.pluginsFn()
}

func (bc *MainFile) Version(ctx context.Context, root string) (string, error) {
	versions := map[string]string{
		"time": time.Now().Format(time.RFC3339),
	}
	m := func() (string, error) {
		b, err := json.Marshal(versions)
		if err != nil {
			return "", plugins.Wrap(bc, err)
		}
		return string(b), nil
	}

	for _, p := range bc.ScopedPlugins() {
		bv, ok := p.(Versioner)
		if !ok {
			continue
		}

		s, err := bv.BuildVersion(ctx, root)
		if err != nil {
			return "", plugins.Wrap(p, err)
		}
		if len(s) == 0 {
			continue
		}
		versions[p.PluginName()] = strings.TrimSpace(s)
	}
	return m()
}

func (bc *MainFile) generateNewMain(ctx context.Context, info here.Info, version string, ws ...io.Writer) error {
	var imports []string
	for _, p := range bc.ScopedPlugins() {
		bi, ok := p.(Importer)
		if !ok {
			continue
		}

		i, err := bi.BuildImports(ctx, info.Dir)
		if err != nil {
			return plugins.Wrap(p, err)
		}
		imports = append(imports, i...)
	}

	if i, err := here.Dir(filepath.Join(info.Dir, "actions")); err == nil {
		imports = append(imports, i.ImportPath)
	}

	sort.Strings(imports)

	bt := struct {
		BuildTime       string
		BuildVersion    string
		Imports         []string
		Info            here.Info
		WithFallthrough bool
	}{
		BuildTime:    strconv.Quote(time.Now().Format(time.RFC3339)),
		BuildVersion: strconv.Quote(version),
		Imports:      imports,
		Info:         info,
	}

	ft := bc.withFallthroughFn
	if ft == nil {
		ft = func() bool {
			c := exec.CommandContext(ctx, "go", "doc", path.Join(info.ImportPath, "cli")+".Buffalo")
			err := c.Run()
			return err == nil
		}
	}
	bt.WithFallthrough = ft()

	t, err := template.New(mainBuildFile).Parse(mainBuildTmpl)
	if err != nil {
		return plugins.Wrap(bc, err)
	}

	f, err := os.Create(filepath.Join(info.Dir, mainBuildFile))
	if err != nil {
		return plugins.Wrap(bc, err)
	}
	defer f.Close()

	if err := t.Execute(io.MultiWriter(append(ws, f)...), bt); err != nil {
		return plugins.Wrap(bc, err)
	}
	return nil
}

func (bc *MainFile) BeforeBuild(ctx context.Context, root string, args []string) error {
	info, err := bc.binaryFolderInfo(root)
	if err != nil {
		return plugins.Wrap(bc, err)
	}

	err = bc.renameMain(info, "main", "originalMain")
	if err != nil {
		return plugins.Wrap(bc, err)
	}

	version, err := bc.Version(ctx, info.Dir)
	if err != nil {
		return plugins.Wrap(bc, err)
	}

	if err := bc.generateNewMain(ctx, info, version); err != nil {
		return plugins.Wrap(bc, err)
	}
	return nil
}

var _ AfterBuilder = &MainFile{}

func (bc *MainFile) AfterBuild(ctx context.Context, root string, args []string, err error) error {
	info, berr := bc.binaryFolderInfo(root)
	if err != nil {
		return plugins.Wrap(bc, berr)
	}

	os.RemoveAll(filepath.Join(info.Dir, mainBuildFile))
	err = bc.renameMain(info, "originalMain", "main")
	return plugins.Wrap(bc, err)
}

func (bc *MainFile) binaryFolderInfo(root string) (here.Info, error) {
	rinfo, err := here.Dir(root)
	if err != nil {
		return here.Info{}, err
	}

	info, err := here.Dir(filepath.Join("cmd", strings.ToLower(rinfo.Name)))
	return info, err
}

func (bc *MainFile) renameMain(info here.Info, from string, to string) error {
	if !info.Module.Main {
		err := fmt.Errorf("module %s is not a main", info.Name)
		return plugins.Wrap(bc, err)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, info.Dir, nil, 0)
	if err != nil {
		return plugins.Wrap(bc, err)
	}

	for _, p := range pkgs {
		for x, f := range p.Files {

			err := func() error {
				var reprint bool
				pre := func(c *astutil.Cursor) bool {
					n := c.Name()
					if n != "Decls" {
						return true
					}

					fd, ok := c.Node().(*ast.FuncDecl)
					if !ok {
						return true
					}

					n = fd.Name.Name
					if n != from {
						return true
					}

					fd.Name = ast.NewIdent(to)
					c.Replace(fd)
					reprint = true
					return true
				}

				res := astutil.Apply(f, pre, nil)
				if !reprint {
					return nil
				}

				f, err := os.Create(x)
				if err != nil {
					return plugins.Wrap(bc, err)
				}
				defer f.Close()
				err = printer.Fprint(f, fset, res)
				return plugins.Wrap(bc, err)
			}()

			if err != nil {
				return plugins.Wrap(bc, err)
			}
		}
	}
	return nil
}
