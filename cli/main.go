package cli

import (
	"context"

	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"
	"github.com/gobuffalo/plugins/plugio"
	"github.com/gobuffalo/plugins/plugprint"
	"github.com/spf13/pflag"
)

// Main is the entry point of the `buffalo` command
func (b *Buffalo) Main(ctx context.Context, root string, args []string) error {
	var help bool
	flags := pflag.NewFlagSet(b.String(), pflag.ContinueOnError)
	flags.BoolVarP(&help, "help", "h", false, "print this help")
	flags.Parse(args)

	pfn := func() []plugins.Plugin {
		return b.Plugins
	}

	plugs := b.Plugins
	for _, p := range plugs {
		if t, ok := p.(Needer); ok {
			t.WithPlugins(pfn)
		}

		if t, ok := p.(StdinNeeder); ok {
			if err := t.SetStdin(plugio.Stdin(plugs...)); err != nil {
				return err
			}
		}

		if t, ok := p.(StdoutNeeder); ok {
			if err := t.SetStdout(plugio.Stdout(plugs...)); err != nil {
				return err
			}
		}

		if t, ok := p.(StderrNeeder); ok {
			if err := t.SetStderr(plugio.Stderr(plugs...)); err != nil {
				return err
			}
		}
	}

	c := plugcmd.FindFromArgs(args, plugs)
	if c != nil {
		return c.Main(ctx, root, args[1:])
	}

	return plugprint.Print(plugio.Stdout(plugs...), b)

}
