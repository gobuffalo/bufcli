package testcmd

import (
	"context"
	"os/exec"

	"github.com/gobuffalo/buffalo-cli/internal/plugins"
)

// Tester is a sub-command of buffalo test.
// 	buffalo test assets
type Tester interface {
	plugins.Plugin
	Test(ctx context.Context, args []string) error
}

type BeforeTester interface {
	plugins.Plugin
	BeforeTest(ctx context.Context, args []string) error
}

type AfterTester interface {
	plugins.Plugin
	AfterTest(ctx context.Context, args []string, err error) error
}

type Runner interface {
	plugins.Plugin
	RunTests(ctx context.Context, cmd *exec.Cmd) error
}

type Tagger interface {
	TestTags(ctx context.Context, root string) ([]string, error)
}
