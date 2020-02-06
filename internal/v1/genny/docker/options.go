package docker

import (
	"fmt"

	bufcli "github.com/gobuffalo/buffalo-cli/v2"
	"github.com/gobuffalo/meta"
)

type Options struct {
	App     meta.App `json:"app"`
	Version string   `json:"version"`
	Style   string   `json:"style"`
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	if len(opts.Version) == 0 {
		opts.Version = bufcli.Version
	}
	if len(opts.Style) == 0 {
		opts.Style = "multi"
	}

	switch opts.Style {
	case "multi", "standard":
	default:
		return fmt.Errorf("unknown style option %s", opts.Style)
	}

	return nil
}
