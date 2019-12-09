package plugins

import "io"

type IO interface {
	StderrGetter
	StderrSetter
	StdinGetter
	StdinSetter
	StdoutGetter
	StdoutSetter
}

type IOGetters interface {
	StderrGetter
	StdinGetter
	StdoutGetter
}

type IOSetters interface {
	StderrSetter
	StdinSetter
	StdoutSetter
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

type StdinSetter interface {
	SetStdin(r io.Reader)
}

type StdoutSetter interface {
	SetStdout(w io.Writer)
}

type StderrSetter interface {
	SetStderr(w io.Writer)
}
