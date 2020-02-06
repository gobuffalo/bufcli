package test

import (
	"context"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/v2/plugins"
)

// Tester is a sub-command of buffalo test.
// 	buffalo test webpack
type Tester interface {
	plugins.Plugin
	Test(ctx context.Context, root string, args []string) error
}

type BeforeTester interface {
	plugins.Plugin
	BeforeTest(ctx context.Context, root string, args []string) error
}

type AfterTester interface {
	plugins.Plugin
	AfterTest(ctx context.Context, root string, args []string, err error) error
}

type Runner interface {
	plugins.Plugin
	RunTests(ctx context.Context, root string, cmd *exec.Cmd) error
}

type Argumenter interface {
	TestArgs(ctx context.Context, root string) ([]string, error)
}
