package plugins

import (
	"context"
	"io"
)

type PluginScoper interface {
	ScopedPlugins() []Plugin
}

type PluginFeeder func() []Plugin

type PluginNeeder interface {
	WithPlugins(PluginFeeder)
}

type Hider interface {
	HidePlugin()
}

type IO interface {
	StderrGetter
	StdinGetter
	StdoutGetter
}

type StdinGetter interface {
	Stdin() io.Reader
}

type StdoutGetter interface {
	Stdout() io.Writer
}

type StderrGetter interface {
	Stderr() io.Writer
}

type Aliases interface {
	Plugin
	Aliases() []string
}

type NamedCommand interface {
	Plugin
	CmdName() string
}

type NamedWriter interface {
	Plugin
	NamedWriter(ctx context.Context, n string) (io.Writer, error)
}
