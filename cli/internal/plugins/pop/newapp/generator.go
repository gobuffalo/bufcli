package newapp

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/markbates/pkger"
	"github.com/spf13/pflag"
)

type Generator struct {
	databaseType string
	skip         bool

	flags *pflag.FlagSet
}

func (g Generator) PluginName() string {
	return "pop/db"
}

func (g Generator) Description() string {
	return "Generates Pop needed files when application is created"
}

func (g *Generator) Newapp(ctx context.Context, root string, name string, args []string) error {
	if g.skip {
		return nil
	}

	err := g.addModels(root, name)
	if err != nil {
		return err
	}

	err = g.addDatabaseConfig(root, name)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) addDatabaseConfig(root, name string) error {
	td := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/internal/plugins/pop/newapp/templates")

	dbtype := g.databaseType
	if dbtype == "" {
		dbtype = "postgres"
	}

	tfile := fmt.Sprintf("database.%v.yml.tmpl", dbtype)
	tf, err := pkger.Open(filepath.Join(td, tfile))
	if err != nil {
		return err
	}

	t, err := ioutil.ReadAll(tf)
	if err != nil {
		return err
	}

	template, err := template.New("database.yml").Parse(string(t))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(root, "database.yml"))
	if err != nil {
		return err
	}

	data := struct {
		Root string
		Name string
	}{
		Root: root + string(filepath.Separator),
		Name: name,
	}

	err = template.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) addModels(root, name string) error {
	td := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/internal/plugins/pop/newapp/templates")
	err := os.Mkdir(filepath.Join(root, "models"), 0777)
	if err != nil {
		return err
	}

	tf, err := pkger.Open(filepath.Join(td, "models.go.tmpl"))
	if err != nil {
		return err
	}

	t, err := ioutil.ReadAll(tf)
	if err != nil {
		return err
	}

	template, err := template.New("models.go").Parse(string(t))
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(root, "models", "models.go"))
	if err != nil {
		return err
	}

	err = template.Execute(f, nil)
	if err != nil {
		return err
	}

	return nil
}
