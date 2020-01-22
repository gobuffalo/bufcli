package tester

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/test"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/here"
	"github.com/gobuffalo/pop/v5"
)

var _ plugins.Plugin = &Tester{}
var _ test.Argumenter = &Tester{}
var _ test.BeforeTester = &Tester{}

type Tester struct {
	info here.Info
}

func (t *Tester) TestArgs(ctx context.Context, root string) ([]string, error) {
	args := []string{"-p", "1"}

	dy := filepath.Join(root, "database.yml")
	if _, err := os.Stat(dy); err != nil {
		return args, nil
	}

	b, err := ioutil.ReadFile(dy)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(b, []byte("sqlite")) {
		args = append(args, "-tags", "sqlite")
	}
	return args, nil
}

func (t *Tester) WithHereInfo(i here.Info) {
	t.info = i
}

func (t *Tester) HereInfo() (here.Info, error) {
	if t.info.IsZero() {
		return here.Current()
	}
	return t.info, nil
}

func (Tester) Name() string {
	return "pop/tester"
}

func (t *Tester) BeforeTest(ctx context.Context, args []string) error {
	info, err := t.HereInfo()
	if err != nil {
		return err
	}

	if err := pop.AddLookupPaths(info.Dir, info.Module.Dir); err != nil {
		return err
	}

	db, ok := ctx.Value("tx").(*pop.Connection)
	if !ok {
		if _, err := os.Stat(filepath.Join(info.Dir, "database.yml")); err != nil {
			return err
		}

		db, err = pop.Connect("test")
		if err != nil {
			return err
		}
	}
	// drop the test db:
	db.Dialect.DropDB()

	// create the test db:
	if err = db.Dialect.CreateDB(); err != nil {
		return err
	}

	for _, a := range args {
		if a == "--force-migrations" {
			return t.forceMigrations(db)
		}
	}

	schema, err := t.findSchema()
	if err != nil {
		return err
	}
	if schema == nil {
		return t.forceMigrations(db)
	}

	if err = db.Dialect.LoadSchema(schema); err != nil {
		return err
	}
	return nil
}

func (t *Tester) forceMigrations(db *pop.Connection) error {
	info, err := t.HereInfo()
	if err != nil {
		return err
	}

	ms := filepath.Join(info.Dir, "migrations")
	fm, err := pop.NewFileMigrator(ms, db)
	if err != nil {
		return err
	}
	return fm.Up()
}

func (t *Tester) findSchema() (io.Reader, error) {
	info, err := t.HereInfo()
	if err != nil {
		return nil, err
	}

	ms := filepath.Join(info.Dir, "migrations", "schema.sql")
	if f, err := os.Open(ms); err == nil {
		return f, nil
	}

	if dev, err := pop.Connect("development"); err == nil {
		schema := &bytes.Buffer{}
		if err = dev.Dialect.DumpSchema(schema); err == nil {
			return schema, nil
		}
	}

	return nil, nil
}
