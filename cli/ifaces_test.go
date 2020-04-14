package cli

import (
	"io"
	"testing"

	"github.com/gobuffalo/plugins"
	"github.com/stretchr/testify/require"
)

type TestPlugin struct{}

func (t TestPlugin) PluginName() string {
	return "test/plugin"
}

func Feeder_Test() []plugins.Plugin {
	return []plugins.Plugin{
		TestPlugin{},
	}
}

type TestNeeder struct {
	Testy               *testing.T
	ExpectedPluginCount int
	Error               error
}

func (tn TestNeeder) PluginName() string {
	return "test/plugin/needer"
}

func (tn TestNeeder) WithPlugins(feeder plugins.Feeder) {
	r := require.New(tn.Testy)
	ps := feeder()
	r.Equal(len(ps), tn.ExpectedPluginCount)
}

func (tn TestNeeder) SetStdin(io.Reader) error {
	return tn.Error
}

func Test_TestNeeder(t *testing.T) {

	tn := TestNeeder{
		Testy:               t,
		ExpectedPluginCount: 1,
	}

	tn.WithPlugins(Feeder_Test)
}

type TestInNeeder struct {
	Error error
}

func (tin TestInNeeder) PluginName() string {
	return "test/plugin/inneeder"
}

func (tin TestInNeeder) SetStdin(io.Reader) error {
	return tin.Error
}

type TestOutNeeder struct {
	Error error
}

func (ton TestOutNeeder) PluginName() string {
	return "test/plugin/outneeder"
}

func (ton TestOutNeeder) SetStdout(io.Writer) error {
	return ton.Error
}

type TestErrNeeder struct {
	Error error
}

func (ten TestErrNeeder) PluginName() string {
	return "test/plugin/errneeder"
}

func (ten TestErrNeeder) SetStderr(io.Writer) error {
	return ten.Error
}
