package buildcmd

import (
	"io"

	"github.com/spf13/pflag"
)

func (bc *BuildCmd) PrintFlags(w io.Writer) error {
	flags := bc.flagSet()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}

func (bc *BuildCmd) flagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet(bc.String(), pflag.ContinueOnError)
	flags.SetOutput(bc.Stdout())

	flags.BoolVar(&bc.skipTemplateValidation, "skip-template-validation", false, "skip validating templates")
	flags.BoolVarP(&bc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&bc.verbose, "verbose", "v", false, "print debugging information")
	flags.BoolVarP(&bc.Static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")

	flags.StringVar(&bc.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	flags.StringVar(&bc.Mod, "mod", "", "-mod flag for go build")
	flags.StringVarP(&bc.Bin, "output", "o", bc.Bin, "set the name of the binary")
	flags.StringVarP(&bc.Environment, "environment", "", "development", "set the environment for the binary")
	flags.StringVarP(&bc.Tags, "tags", "t", "", "compile with specific build tags")

	plugs := bc.Plugins()

	for _, p := range plugs {
		switch t := p.(type) {
		case BuildFlagger:
			for _, f := range t.BuildFlags() {
				flags.AddGoFlag(f)
			}
		case BuildPflagger:
			for _, f := range t.BuildFlags() {
				flags.AddFlag(f)
			}
		}
	}

	return flags
}
