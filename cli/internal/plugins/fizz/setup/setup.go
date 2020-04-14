package setup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/setup"
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/pop/v5"
)

var _ plugins.Plugin = &Setup{}
var _ setup.Migrater = &Setup{}

type Setup struct{}

func (s *Setup) PluginName() string {
	return "fizz/setup"
}

func (s *Setup) MigrateDB(ctx context.Context, root string, args []string) error {
	if err := pop.LoadConfigFile(); err != nil {
		return err
	}

	env := os.Getenv("GO_ENV")
	if len(env) == 0 {
		env = "development"
	}

	conn, ok := pop.Connections[env]
	if !ok {
		return fmt.Errorf("no connection found for %s", env)
	}

	mg := filepath.Join(root, "migrations")
	mig, err := pop.NewFileMigrator(mg, conn)
	if err != nil {
		return err
	}
	return mig.Up()
}
