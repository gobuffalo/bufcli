package bzr

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-cli/v2/cli/internal/plugins/build"
	"github.com/gobuffalo/buffalo-cli/v2/plugins"
	"github.com/gobuffalo/buffalo-cli/v2/plugins/plugprint"
)

type Buffalo struct{}

var _ build.Versioner = &Buffalo{}

func (b *Buffalo) BuildVersion(ctx context.Context, root string) (string, error) {
	if _, err := exec.LookPath("bzr"); err != nil {
		return "", err
	}

	bb := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, "bzr", "revno")
	cmd.Stdout = bb
	cmd.Stderr = bb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %s", err, bb.String())
	}
	s := strings.TrimSpace(bb.String())
	if len(s) == 0 {
		return "", nil
	}
	return s, nil
}

var _ plugins.Plugin = Buffalo{}

// Name is the name of the plugin.
// This will also be used for the cli sub-command
// 	"pop" | "heroku" | "auth" | etc...
func (b Buffalo) Name() string {
	return "bzr"
}

var _ plugprint.Describer = Buffalo{}

func (b Buffalo) Description() string {
	return "Provides bzr related hooks to Buffalo applications."
}
