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

func (d *Daggerverse) Release(
	ctx context.Context,
	// +optional
	// +default=false
	dryRun bool,
) ([]string, error) {
	return dag.DaggerverseCockpit().Publish(ctx, d.Source, dagger.DaggerverseCockpitPublishOpts{
		DryRun: dryRun,
		Exclude: []string{
			"dagger.json", // don't include our own CI
			"artifactory/examples/go",
			"jfrogcli/examples/go",
		},
	})
}
