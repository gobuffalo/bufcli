package built

import (
	"context"

	"github.com/gobuffalo/plugins"
)

var _ plugins.Plugin = Initer{}

type Initer struct{}

func (Initer) PluginName() string {
	return "pop/built/initer"
}

func (p *Initer) BuiltInit(ctx context.Context, root string, args []string) error {
	return nil
}
