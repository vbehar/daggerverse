package main

import (
	"context"

	"github.com/vbehar/daggerverse/gitlab-cli/internal/dagger"
)

// Release allows you to interact with GitLab Releases.
func (g *GitlabCli) Release(
	// name of the tag.
	tagName string,
	// a file from which to read the release description.
	// +optional
	descriptionFile *dagger.File,
	// container to use for the command, instead of the default container
	// you can use this to customize the container
	// +optional
	ctr *dagger.Container,
) *Release {
	return &Release{
		GitlabCli:       g,
		TagName:         tagName,
		DescriptionFile: descriptionFile,
		Ctr:             ctr,
	}
}

// Release allows you to interact with GitLab Releases.
type Release struct {
	// +private
	GitlabCli       *GitlabCli
	TagName         string
	DescriptionFile *dagger.File
	Ctr             *dagger.Container
}

func (r *Release) descriptionFileName() string {
	if r.DescriptionFile != nil {
		return "description.md"
	}
	return ""
}

// Container returns a container ready to be used for managing releases.
func (r *Release) Container(ctx context.Context) *dagger.Container {
	ctr := r.Ctr
	if ctr == nil {
		ctr = r.GitlabCli.Container(ctx)
	}

	if r.DescriptionFile != nil {
		ctr = ctr.
			WithWorkdir("/workdir").
			WithMountedFile(r.descriptionFileName(), r.DescriptionFile)
	}

	return ctr
}

// CreateRelease creates a new release for the given tag.
// If the tag doesn't exist, it will be created - if you also provide a gitRef.
func (r *Release) Create(
	ctx context.Context,
	// if the tag should be created, it will be created from this ref.
	// can be a commit or a branch.
	// +optional
	gitRef string,
) (string, error) {
	return r.Container(ctx).
		WithExec([]string{
			"release-cli", "create",
			"--tag-name", r.TagName,
			"--description", r.descriptionFileName(),
			"--ref", gitRef,
		}).
		Stdout(ctx)
}

// UpdateRelease updates an existing release for the given tag.
func (r *Release) Update(
	ctx context.Context,
) (string, error) {
	return r.Container(ctx).
		WithExec([]string{
			"release-cli", "update",
			"--tag-name", r.TagName,
			"--description", r.descriptionFileName(),
		}).
		Stdout(ctx)
}
