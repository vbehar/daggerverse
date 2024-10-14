package main

import (
	"context"

	"github.com/vbehar/daggerverse/jfrogcli/examples/go/internal/dagger"
)

type Examples struct{}

func (e *Examples) JFrogCLI_Install() *dagger.Container {
	return dag.Jfrogcli().Install()
}

func (e *Examples) JFrogCLI_InstallVersion(version string) *dagger.Container {
	return dag.Jfrogcli(dagger.JfrogcliOpts{
		Version: version,
	}).Install()
}

func (e *Examples) JFrogCLI_InstallInto(ctr *dagger.Container) *dagger.Container {
	return dag.Jfrogcli().Install(dagger.JfrogcliInstallOpts{
		Base: ctr,
	})
}

func (e *Examples) JFrogCLI_Run(ctx context.Context) (string, error) {
	return dag.Jfrogcli().Install().
		WithExec([]string{"jf", "--version"}).
		Stdout(ctx)
}
