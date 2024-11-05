package main

import (
	"context"

	"github.com/vbehar/daggerverse/git-info/examples/go/internal/dagger"
)

type Examples struct{}

func (e *Examples) GitInfo_Version(
	ctx context.Context,
	gitDirectory *dagger.Directory,
) {
	version, err := dag.GitInfo(gitDirectory).Version(ctx)
	if err != nil {
		panic(err)
	}

	println("current git version:", version)
}

func (e *Examples) GitInfo_SetEnvVariablesOnContainer(
	gitDirectory *dagger.Directory,
) {
	ctr := dag.Container().From("alpine:latest")

	ctr = dag.GitInfo(gitDirectory).SetEnvVariablesOnContainer(ctr)

	ctr.WithExec([]string{
		"sh", "-c",
		"echo do something with $GIT_TAG or other env vars...",
	})
}
