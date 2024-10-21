// A Dagger Module to install and run the GitLab CLI.
//
// GitLab CLI is a command-line tool that allows you to interact with the GitLab API.
// See https://gitlab.com/gitlab-org/cli for more information.
// This module also contains the (deprecated) GitLab Release CLI.
// See https://gitlab.com/gitlab-org/release-cli for more information.
package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/vbehar/daggerverse/gitlab-cli/internal/dagger"
)

// GitlabCli is a Dagger Module to interact with the GitLab CLI.
type GitlabCli struct {
	PrivateToken      *dagger.Secret
	JobToken          *dagger.Secret
	Host              string
	Repo              string
	Group             string
	ReleaseCliVersion string
}

func New(
	// private (personal) token to use for authentication with GitLab.
	// +optional
	privateToken *dagger.Secret,
	// (CI) job token to use for authentication with GitLab.
	// Defined as CI_JOB_TOKEN in the GitLab CI environment.
	// +optional
	jobToken *dagger.Secret,
	// host of the GitLab instance.
	// +optional
	host string,
	// default gitlab repository for commands accepting the --repo flag.
	// +optional
	repo string,
	// default gitlab group for commands accepting the --group flag.
	// +optional
	group string,
	// version of the GitLab Release CLI tool to use.
	// https://gitlab.com/gitlab-org/release-cli
	// +optional
	// +default="v0.18.0"
	releaseCliVersion string,
) *GitlabCli {
	return &GitlabCli{
		PrivateToken:      privateToken,
		JobToken:          jobToken,
		Host:              host,
		Repo:              repo,
		Group:             group,
		ReleaseCliVersion: releaseCliVersion,
	}
}

// Container returns a container with the GitLab CLI installed.
func (g *GitlabCli) Container(
	ctx context.Context,
) *dagger.Container {
	ctr := dag.Container().
		From("cgr.dev/chainguard/wolfi-base").
		WithExec([]string{"apk", "add", "--update", "--no-cache",
			"ca-certificates",
			"glab",
			"jq",
		}).
		WithFile("/usr/bin/release-cli", g.releaseCLI(ctx)).
		WithExec([]string{"chmod", "+x", "/usr/bin/release-cli"})

	if g.PrivateToken != nil {
		ctr = ctr.
			WithSecretVariable("GITLAB_TOKEN", g.PrivateToken).        // for glab
			WithSecretVariable("GITLAB_PRIVATE_TOKEN", g.PrivateToken) // for release-cli
	}
	if g.JobToken != nil {
		ctr = ctr.WithSecretVariable("CI_JOB_TOKEN", g.JobToken) // for release-cli
	}
	if g.Host != "" {
		ctr = ctr.
			WithEnvVariable("GITLAB_HOST", g.Host).  // for glab
			WithEnvVariable("CI_SERVER_URL", g.Host) // for release-cli
	}
	if g.Repo != "" {
		ctr = ctr.
			WithEnvVariable("GITLAB_REPO", g.Repo).                  // for glab
			WithEnvVariable("CI_PROJECT_ID", url.PathEscape(g.Repo)) // for release-cli
	}
	if g.Group != "" {
		ctr = ctr.WithEnvVariable("GITLAB_GROUP", g.Group) // for glab
	}

	return ctr
}

// Run runs the glab CLI with the given arguments.
func (g *GitlabCli) Run(
	ctx context.Context,
	// arguments to pass to the glab CLI
	// +optional
	args []string,
	// container to use for the command, instead of the default container
	// you can use this to customize the container
	// +optional
	ctr *dagger.Container,
) (string, error) {
	if ctr == nil {
		ctr = g.Container(ctx)
	}

	return ctr.
		WithEntrypoint([]string{"glab"}).
		WithExec(args, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		Stdout(ctx)
}

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
			WithFile(r.descriptionFileName(), r.DescriptionFile)
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
		WithFocus().
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
		WithFocus().
		WithExec([]string{
			"release-cli", "update",
			"--tag-name", r.TagName,
			"--description", r.descriptionFileName(),
		}).
		Stdout(ctx)
}

func (g *GitlabCli) releaseCLI(ctx context.Context) *dagger.File {
	platform, err := dag.DefaultPlatform(ctx)
	if err != nil {
		platform = dagger.Platform("linux/amd64")
	}
	var os, arch string
	elems := strings.Split(string(platform), "/")
	if len(elems) >= 2 {
		os = elems[0]
		arch = elems[1]
	} else {
		os = "linux"
		arch = "amd64"
	}

	return dag.HTTP(fmt.Sprintf(
		"https://gitlab.com/gitlab-org/release-cli/-/releases/%s/downloads/bin/release-cli-%s-%s",
		g.ReleaseCliVersion, os, arch,
	))
}
