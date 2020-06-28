package newapp

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AddModels(t *testing.T) {
	r := require.New(t)

	root, err := ioutil.TempDir("", "")
	r.NoError(err)

	g := Generator{}
	g.addModels(root, "application")

	r.DirExists(filepath.Join(root, "models"))
	r.FileExists(filepath.Join(root, "models", "models.go"))

	f, err := ioutil.ReadFile(filepath.Join(filepath.Join(root, "models", "models.go")))
	r.NoError(err)

	r.Contains(string(f), "// DB is a connection to your database to be used")

}

func Test_addDatabaseConfig(t *testing.T) {
	r := require.New(t)
	root, err := ioutil.TempDir("", "")
	r.NoError(err)

	tcases := []struct {
		databaseType    string
		expectedContent []string
	}{
		{"", []string{"dialect: postgres"}},
		{"postgres", []string{"dialect: postgres"}},
		{"mariadb", []string{`dialect: "mariadb"`}},
		{"mysql", []string{`dialect: "mysql"`}},
		{"cockroach", []string{`dialect: cockroach`}},
		{"sqlite3", []string{`dialect: "sqlite3"`, filepath.Join(root, "application_development.sqlite")}},
	}

	for _, tcase := range tcases {
		g := Generator{
			databaseType: tcase.databaseType,
			skip:         false,
		}

		err = g.addDatabaseConfig(root, "application")
		r.NoError(err, tcase.databaseType)

		r.FileExists(filepath.Join(root, "database.yml"), tcase.databaseType)
		f, err := ioutil.ReadFile(filepath.Join(filepath.Join(root, "database.yml")))
		r.NoError(err, tcase.databaseType)

		for _, c := range tcase.expectedContent {
			r.Contains(string(f), c, tcase.databaseType)
		}
	}

}

func Test_Newapp(t *testing.T) {
	r := require.New(t)

	root, err := ioutil.TempDir("", "")
	r.NoError(err)

	g := Generator{}

	err = g.Newapp(context.Background(), root, "application", []string{})
	r.NoError(err)

	r.DirExists(filepath.Join(root, "models"))
	r.FileExists(filepath.Join(root, "models", "models.go"))
	r.FileExists(filepath.Join(root, "database.yml"))

	f, err := ioutil.ReadFile(filepath.Join(filepath.Join(root, "database.yml")))
	r.NoError(err)

	r.Contains(string(f), "dialect: postgres")

	g.skip = true
	root, err = ioutil.TempDir("", "")
	r.NoError(err)

	err = g.Newapp(context.Background(), root, "application", []string{})
	r.NoError(err)

	r.NoDirExists(filepath.Join(root, "models"))
	r.NoFileExists(filepath.Join(root, "models", "models.go"))
	r.NoFileExists(filepath.Join(root, "database.yml"))

}
