package pop

import (
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/build"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/built"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/generators/actions"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/generators/actiontest"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/generators/models"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/newapp"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/setup"
	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/pop/test"
	"github.com/gobuffalo/plugins"
)

func Plugins() []plugins.Plugin {
	return []plugins.Plugin{
		&actions.Generator{},
		&actiontest.Generator{},
		&build.Builder{},
		&built.Initer{},
		&models.Generator{},
		&setup.Setup{},
		&test.Tester{},
		&newapp.Generator{},
	}
}
