package newapp

import (
	"io"
	"io/ioutil"

	"github.com/gobuffalo/buffalo-cli/v2/internal/flagger"
	"github.com/spf13/pflag"
)

func (g *Generator) NewappFlags() []*pflag.Flag {
	return flagger.SetToSlice(g.Flags())
}

func (g *Generator) Flags() *pflag.FlagSet {
	if g.flags != nil {
		return g.flags
	}

	flg := pflag.NewFlagSet(g.PluginName(), pflag.ContinueOnError)
	flg.SetOutput(ioutil.Discard)

	flg.StringVarP(&g.databaseType, "type", "", "postgres", "specify the type of database you want to use [cockroach, mariadb, mysql, postgres]")
	flg.BoolVarP(&g.skip, "skip-pop", "", false, "skips adding pop/soda to your app")

	g.flags = flg
	return g.flags
}

func (g *Generator) PrintFlags(w io.Writer) error {
	flags := g.Flags()
	flags.SetOutput(w)
	flags.PrintDefaults()
	return nil
}
