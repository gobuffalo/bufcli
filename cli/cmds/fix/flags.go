package fix

import (
	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/gobuffalo/plugins/plugflag"
	"github.com/spf13/pflag"
)

func (cmd *Cmd) Flags() *pflag.FlagSet {
	if cmd.flags != nil {
		return cmd.flags
	}

	flags := pflag.NewFlagSet(cmd.PluginName(), pflag.ContinueOnError)
	flags.BoolVarP(&cmd.help, "help", "h", false, "print this help")

	for _, p := range cmd.ScopedPlugins() {
		switch t := p.(type) {
		case Flagger:
			for _, f := range plugflag.Clean(p, t.FixFlags()) {
				flags.AddGoFlag(f)
			}
		case Pflagger:
			for _, f := range flagger.CleanPflags(p, t.FixFlags()) {
				flags.AddGoFlag(f)
			}
		}
	}

	cmd.flags = flags

	return cmd.flags
}
