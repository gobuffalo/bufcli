package actiontest

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/meta/v2"
)

func (mg *Generator) GenerateResourceActionTests(ctx context.Context, root string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must specify a resource")
	}

	info, err := mg.HereInfo()
	if err != nil {
		return err
	}

	resourceName := args[0]

	testPkg := mg.TestPkg
	if len(testPkg) == 0 {
		testPkg = "actions"
	}

	app, err := meta.New(info)
	if err != nil {
		return err
	}

	actions := []name.Ident{
		name.New("list"),
		name.New("show"),
		name.New("create"),
		name.New("update"),
		name.New("destroy"),
	}

	if app.As["web"] {
		actions = append(actions, name.New("new"), name.New("edit"))
	}

	pres := struct {
		Actions []name.Ident
		Name    name.Ident
		TestPkg string
	}{
		Actions: actions,
		Name:    name.New(resourceName),
		TestPkg: testPkg,
	}

	t, err := template.New(resourceName).Parse(actionsTestTmpl)
	if err != nil {
		return err
	}

	p := fmt.Sprintf("%s_test.go", pres.Name.Folder().Pluralize())
	fp := filepath.Join(root, "actions", p)

	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return err
	}

	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(f, pres); err != nil {
		return err
	}

	return nil
}
