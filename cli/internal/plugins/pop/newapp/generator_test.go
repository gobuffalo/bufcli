package newapp

import (
	"context"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/markbates/pkger"
)

type Generator struct {
	databaseType string
}

func (g Generator) PluginName() string {
	return "pop/db"
}

func (g Generator) Description() string {
	return "Generates Pop needed files when application is created"
}

func (mg *Generator) Newapp(ctx context.Context, root string, name string, args []string) error {
	err := mg.addModels(root, name)
	if err != nil {
		return err
	}

	err = mg.addDatabaseConfig(root, name)
	if err != nil {
		return err
	}

	return nil
}

func (mg *Generator) addDatabaseConfig(root, name string) error {
	td := pkger.Include("github.com/gobuffalo/buffalo-cli/v2:/cli/internal/plugins/pop/newapp/templates")

	tf, err := pkger.Open(filepath.Join(td, "database.postgres.yml.tmpl"))
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
		Name   string
		Prefix string
	}{
		Name:   name,
		Prefix: root,
	}

	err = template.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

func (mg *Generator) addModels(root, name string) error {
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
