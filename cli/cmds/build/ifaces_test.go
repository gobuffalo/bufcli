package build

import (
	"context"
	"flag"
	"os/exec"

	"github.com/spf13/pflag"
)

var _ Builder = &builder{}

type builder struct {
	name string
	args []string
	err  error
}

func (b *builder) PluginName() string {
	if len(b.name) == 0 {
		return "builder"
	}
	return b.name
}

func (b *builder) Build(ctx context.Context, root string, args []string) error {
	b.args = args
	return b.err
}

var _ BeforeBuilder = &beforeBuilder{}

type beforeBuilder struct {
	name string
	args []string
	err  error
}

func (b *beforeBuilder) PluginName() string {
	if len(b.name) == 0 {
		return "beforeBuilder"
	}
	return b.name
}

func (b *beforeBuilder) BeforeBuild(ctx context.Context, root string, args []string) error {
	b.args = args
	return b.err
}

var _ AfterBuilder = &afterBuilder{}

type afterBuilder struct {
	name string
	args []string
	err  error
}

func (b *afterBuilder) PluginName() string {
	if len(b.name) == 0 {
		return "afterBuilder"
	}
	return b.name
}

func (b *afterBuilder) AfterBuild(ctx context.Context, root string, args []string, err error) error {
	b.args = args
	b.err = err
	return err
}

var _ Flagger = &buildFlagger{}

type buildFlagger struct {
	name  string
	flags []*flag.Flag
}

func (b *buildFlagger) PluginName() string {
	if len(b.name) == 0 {
		return "buildFlagger"
	}
	return b.name
}

func (b *buildFlagger) BuildFlags() []*flag.Flag {
	return b.flags
}

var _ Pflagger = &buildPflagger{}

type buildPflagger struct {
	name  string
	flags []*pflag.Flag
}

func (b *buildPflagger) PluginName() string {
	if len(b.name) == 0 {
		return "buildPflagger"
	}
	return b.name
}

func (b *buildPflagger) BuildFlags() []*pflag.Flag {
	return b.flags
}

var _ Packager = &packager{}

type packager struct {
	name  string
	root  string
	files []string
	err   error
}

func (b *packager) PluginName() string {
	if len(b.name) == 0 {
		return "packager"
	}
	return b.name
}

func (b *packager) Package(ctx context.Context, root string, files []string) error {
	b.root = root
	b.files = append(b.files, files...)
	return b.err
}

var _ PackFiler = &packFiler{}

type packFiler struct {
	name  string
	root  string
	files []string
	err   error
}

func (b *packFiler) PackageFiles(ctx context.Context, root string) ([]string, error) {
	b.root = root
	return b.files, b.err
}

func (b *packFiler) PluginName() string {
	if len(b.name) == 0 {
		return "packFiler"
	}
	return b.name
}

var _ Versioner = &buildVersioner{}

type buildVersioner struct {
	name    string
	version string
	root    string
	err     error
}

func (b *buildVersioner) PluginName() string {
	if len(b.name) == 0 {
		return "buildVersioner"
	}
	return b.name
}

func (b *buildVersioner) BuildVersion(ctx context.Context, root string) (string, error) {
	b.root = root
	return b.version, b.err
}

var _ Importer = &buildImporter{}

type buildImporter struct {
	name    string
	imports []string
	root    string
	err     error
}

func (b *buildImporter) PluginName() string {
	if len(b.name) == 0 {
		return "buildImporter"
	}
	return b.name
}

func (b *buildImporter) BuildImports(ctx context.Context, root string) ([]string, error) {
	b.root = root
	return b.imports, b.err
}

var _ Runner = &bladeRunner{}

type bladeRunner struct {
	cmd *exec.Cmd
	err error
}

func (bladeRunner) PluginName() string {
	return "blade"
}

func (b *bladeRunner) RunBuild(ctx context.Context, cmd *exec.Cmd) error {
	b.cmd = cmd
	return b.err
}
