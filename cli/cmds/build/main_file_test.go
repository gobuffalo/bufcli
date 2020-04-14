package build

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/here"
	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

func Test_MainFile_Version(t *testing.T) {
	r := require.New(t)

	bc := &MainFile{}

	ctx := context.Background()

	s, err := bc.Version(ctx, "")
	r.NoError(err)
	r.Contains(s, `"time":`)

	bc.WithPlugins(func() []plugins.Plugin {
		return plugins.Plugins{
			&buildBuilder{},
			&buildVersioner{version: "v1"},
		}
	})
	bc.HidePlugin()

	s, err = bc.Version(ctx, "")
	r.NoError(err)
	r.Contains(s, `"time":`)
	r.Contains(s, `"buildVersioner":"v1"`)
}

func Test_MainFile_Version_Err(t *testing.T) {
	r := require.New(t)

	bc := &MainFile{}

	ctx := context.Background()

	bc.WithPlugins(func() []plugins.Plugin {
		return plugins.Plugins{
			&buildBuilder{},
			&buildVersioner{version: ""},
			&buildVersioner{err: errors.New("Bad Version")},
		}
	})

	s, err := bc.Version(ctx, "")
	r.Error(err)
	r.Contains(err.Error(), "Bad Version")
	r.Equal("", s)
}

func Test_MainFile_generateNewMain(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "")
	defer os.RemoveAll(filepath.Join(ref.Dir, mainBuildFile))

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	ctx := context.Background()
	bb := &bytes.Buffer{}
	err := bc.generateNewMain(ctx, ref, "v1", bb)
	r.NoError(err)

	out := bb.String()
	r.Contains(out, `appcli "github.com/markbates/coke/cli"`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `appcli.Buffalo`)
	r.Contains(out, `originalMain`)

	bc = &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}
	err = bc.generateNewMain(ctx, ref, "v1", bb)
	r.NoError(err)

	out = bb.String()
	r.Contains(out, `appcli "github.com/markbates/coke/cli"`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `appcli.Buffalo`)
	r.Contains(out, `originalMain`)
}

func Test_MainFile_generateNewMain_noCli(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "")
	defer os.RemoveAll(filepath.Join(ref.Dir, mainBuildFile))

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return false },
	}

	bb := &bytes.Buffer{}
	err := bc.generateNewMain(context.Background(), ref, "v1", bb)
	r.NoError(err)

	out := bb.String()
	r.NotContains(out, `appcli "github.com/markbates/coke/cli"`)
	r.NotContains(out, `appcli.Buffalo`)
	r.Contains(out, `_ "github.com/markbates/coke/actions"`)
	r.Contains(out, `originalMain`)
	r.Contains(out, `cb.Main`)
}

func Test_MainFile_generateNewMain_Errs(t *testing.T) {
	r := require.New(t)

	ref := newRef(t, "")
	defer os.RemoveAll(filepath.Join(ref.Dir, mainBuildFile))

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
			err: errors.New("Bad Import"),
		},
		&buildTagger{},
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	ctx := context.Background()
	bb := &bytes.Buffer{}

	err := bc.generateNewMain(ctx, ref, "v1", bb)
	r.Error(err)
	r.Contains(err.Error(), "Bad Import")

	plugs = plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
		&buildTagger{},
	}
	bc = &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
	}

	//Bad writer
	bw := &badWriter{
		err: errors.New("Bad Writer"),
	}
	err = bc.generateNewMain(ctx, ref, "v1", bw)
	r.Error(err)
	r.Contains(err.Error(), "Bad Writer")
}

func Test_MainFile_BeforeBuild(t *testing.T) {
	r := require.New(t)

	root := filepath.Join(".", "coke")
	current_root, err := os.Getwd()
	r.NoError(err)

	err, ref := setupProject(root, t, beforemaingofile)
	defer teardownProject(root, current_root)
	r.NoError(err)

	plugs := plugins.Plugins{
		&buildImporter{
			imports: []string{
				path.Join(ref.ImportPath, "actions"),
			},
		},
	}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	ctx := context.Background()
	var args []string
	err = bc.BeforeBuild(ctx, root, args)
	r.NoError(err)

	raw, err := ioutil.ReadFile("main.go")
	r.NoError(err)
	r.Contains(string(raw), "func originalMain()")
}

func Test_MainFile_BeforeBuild_Err(t *testing.T) {
	r := require.New(t)

	plugs := plugins.Plugins{}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	//Test not in the folder with main (i.e. this folder)
	ctx := context.Background()
	var args []string
	err := bc.BeforeBuild(ctx, "", args)
	r.Error(err)
	r.Contains(err.Error(), "is not a main")

	//Test sending in a plugin with a bad version
	root := filepath.Join(".", "coke")
	current_root, err := os.Getwd()
	r.NoError(err)

	err, _ = setupProject(root, t, beforemaingofile)
	defer teardownProject(root, current_root)
	r.NoError(err)

	plugs = plugins.Plugins{
		&buildVersioner{err: errors.New("Bad Versioner")},
	}
	bc = &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	err = bc.BeforeBuild(ctx, root, args)
	r.Error(err)
	r.Contains(err.Error(), "Bad Versioner")

	// Test error in generate
	plugs = plugins.Plugins{
		&buildImporter{err: errors.New("Bad Importer")},
	}
	bc = &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	err = bc.BeforeBuild(ctx, root, args)
	r.Error(err)
	r.Contains(err.Error(), "Bad Importer")
}

func Test_MainFile_AfterBuild(t *testing.T) {
	r := require.New(t)

	root := filepath.Join(".", "coke")
	current_root, err := os.Getwd()
	r.NoError(err)

	err, _ = setupProject(root, t, aftermaingofile)
	defer teardownProject(root, current_root)
	r.NoError(err)

	plugs := plugins.Plugins{}
	bc := &MainFile{
		pluginsFn: func() []plugins.Plugin {
			return plugs
		},
		withFallthroughFn: func() bool { return true },
	}

	ctx := context.Background()
	var args []string
	err = bc.AfterBuild(ctx, ".", args, nil)
	r.NoError(err)

	raw, err := ioutil.ReadFile("main.go")
	r.NoError(err)
	r.Contains(string(raw), "func main()")
}

func Test_MainFile_AfterBuild_Err(t *testing.T) {
	r := require.New(t)
	bc := &MainFile{}

	ctx := context.Background()
	var args []string
	err := bc.AfterBuild(ctx, "random1370498nc19c", args, nil)
	r.Error(err)
	if runtime.GOOS == "windows" {
		r.Contains(err.Error(), "The system cannot find the file specified")
	} else {
		r.Contains(err.Error(), "no such file or directory")
	}
}

func Test_MainFile_renameMain_err(t *testing.T) {
	r := require.New(t)

	root := filepath.Join(".", "coke")
	current_root, err := os.Getwd()
	r.NoError(err)

	err, _ = setupProject(root, t, aftermaingofile)
	defer teardownProject(root, current_root)
	r.NoError(err)

	// Shouldn't be able to rename in a read-only file
	r.NoError(os.Chmod("main.go", 0444))
	r.NoError(os.Chdir(current_root))

	bc := &MainFile{}

	info, err := here.Dir(root)
	r.NoError(err)
	err = bc.renameMain(info, "originalMain", "main")
	r.Error(err)
	if runtime.GOOS == "windows" {
		r.Contains(err.Error(), "Access is denied.")
	} else {
		r.Contains(err.Error(), "permission denied")
	}

	raw, err := ioutil.ReadFile(filepath.Join(root, "main.go"))
	r.NoError(err)
	r.Contains(string(raw), "func originalMain()")

	// cover when the thing to that would be renamed isn't actually a function
	r.NoError(os.Chmod(filepath.Join(root, "main.go"), 0660))
	r.NoError(ioutil.WriteFile(filepath.Join(root, "main.go"), []byte(renamemaingofile), 0660))

	info, err = here.Dir(root)
	r.NoError(err)
	err = bc.renameMain(info, "originalMain", "main")
	r.NoError(err)

	raw, err = ioutil.ReadFile(filepath.Join(root, "main.go"))
	r.NoError(err)
	r.Contains(string(raw), "func main()")

}

func setupProject(root string, t *testing.T, maincontent string) (error, here.Info) {
	if err := os.MkdirAll(root, os.ModeDir|os.ModePerm); err != nil {
		return err, here.Info{}
	}
	if err := os.Chdir(root); err != nil {
		return err, here.Info{}
	}

	//create files
	if err := ioutil.WriteFile("main.go", []byte(maincontent), 0660); err != nil {
		return err, here.Info{}
	}
	if err := ioutil.WriteFile("go.mod", []byte(gomodfile), 0660); err != nil {
		return err, here.Info{}
	}
	return nil, newRef(t, root)
}

func teardownProject(root, newRoot string) {
	os.Chdir(newRoot)
	os.RemoveAll(root)
}

const beforemaingofile = `
package main

func main() {
	print("go time!")
}
`

const aftermaingofile = `
package main

func originalMain() {
	print("go time!")
}
`

const renamemaingofile = `
package main

var originalMain string = "oops"
func main(){
	print(originalMain)
}	
`

const gomodfile = `
module github.com/markbates/coke

go 1.13`
