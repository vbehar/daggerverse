package main

import (
	"context"

	"github.com/vbehar/daggerverse/jx-release-version/examples/go/internal/dagger"
)

type Examples struct{}

func (e *Examples) JxReleaseVersion_NextVersion(
	ctx context.Context,
	gitDirectory *dagger.Directory,
) {
	nextVersion, err := dag.JxReleaseVersion().NextVersion(ctx, gitDirectory)
	if err != nil {
		panic(err)
	}

	println("next version:", nextVersion)
}

func (e *Examples) JxReleaseVersion_Tag(
	ctx context.Context,
	gitDirectory *dagger.Directory,
	gitToken *dagger.Secret,
) {
	newTag, err := dag.JxReleaseVersion().Tag(ctx, gitDirectory, gitToken, dagger.JxReleaseVersionTagOpts{
		PushTag: true,
	})
	if err != nil {
		panic(err)
	}

	println("new tag pushed:", newTag)
}
