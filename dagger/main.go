package main

import (
	"context"

	"dagger/daggerverse/internal/dagger"
)

type Daggerverse struct {
	Source *dagger.Directory
}

func New(
	source *dagger.Directory,
) *Daggerverse {
	return &Daggerverse{
		Source: source,
	}
}

func (d *Daggerverse) Tag(
	ctx context.Context,
	gitToken *dagger.Secret,
) (string, error) {
	return dag.JxReleaseVersion().Tag(ctx,
		d.Source.Directory(".git"),
		gitToken,
		dagger.JxReleaseVersionTagOpts{
			FetchTags: true,
		},
	)
}
