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

const (
	// use fixed base images for reproductible builds and improved caching
	// the base image: https://images.chainguard.dev/directory/image/wolfi-base/overview
	// retrieve the latest sha256 hash with: `crane digest cgr.dev/chainguard/wolfi-base:latest`
	// and to retrieve its creation time: `crane config cgr.dev/chainguard/wolfi-base:latest | jq .created`
	// This one is from 2025-05-22T20:15:28Z
	baseWolfiImage = "cgr.dev/chainguard/wolfi-base:latest@sha256:0c35d31660ee8ff26c0893f7f1fe5752aea11f036536368791d2854e67112f85"
)

// GitlabCli is a Dagger Module to interact with the GitLab CLI.
type GitlabCli struct {
	PrivateToken      *dagger.Secret
	JobToken          *dagger.Secret
	Host              string
	Repo              string
	Group             string
	GitDirectory      *dagger.Directory
	ReleaseCliVersion string
	GLabVersion       string
	GLabDebug         bool
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
	// Git directory to use as a context.
	// The GitLab CLI will retrieve the repository, branch, etc from this directory.
	// +optional
	gitDirectory *dagger.Directory,
	// version of the GitLab Release CLI tool to use.
	// https://gitlab.com/gitlab-org/release-cli/-/releases
	// https://gitlab.com/gitlab-org/release-cli/-/blob/master/CHANGELOG.md
	// +optional
	// +default="v0.23.0"
	releaseCliVersion string,
	// version of the GitLab CLI tool to use.
	// https://gitlab.com/gitlab-org/cli/-/releases
	// +optional
	// +default="1.58.0"
	glabVersion string,
	// enable debug mode for the GitLab CLI.
	// +optional
	// +default=false
	glabDebug bool,
) *GitlabCli {
	return &GitlabCli{
		PrivateToken:      privateToken,
		JobToken:          jobToken,
		Host:              host,
		Repo:              repo,
		Group:             group,
		GitDirectory:      gitDirectory,
		ReleaseCliVersion: releaseCliVersion,
		GLabVersion:       glabVersion,
		GLabDebug:         glabDebug,
	}
}

// Container returns a container with the GitLab CLI installed.
func (g *GitlabCli) Container(
	ctx context.Context,
) *dagger.Container {
	ctr := dag.Container().
		From(baseWolfiImage).
		WithExec([]string{"apk", "add", "--update", "--no-cache",
			"ca-certificates",
			"glab~=" + g.GLabVersion,
			"git",
			"jq",
		}).
		WithFile("/usr/bin/release-cli", g.releaseCLI(ctx)).
		WithExec([]string{"chmod", "+x", "/usr/bin/release-cli"})

	if g.GLabDebug {
		ctr = ctr.WithEnvVariable("DEBUG", "true") // for glab
	}
	if g.Host != "" {
		ctr = ctr.
			WithEnvVariable("GITLAB_HOST", hostname(g.Host)).   // for glab
			WithEnvVariable("CI_SERVER_URL", serverURL(g.Host)) // for release-cli
	}
	if g.JobToken != nil {
		ctr = ctr.WithSecretVariable("CI_JOB_TOKEN", g.JobToken) // for release-cli
		// for glab, we need to run the login cmd
		ctr = ctr.
			WithExec([]string{
				"/bin/sh", "-c",
				"glab auth login --hostname " + hostname(g.Host) + " --job-token $CI_JOB_TOKEN",
			})
	}
	if g.PrivateToken != nil {
		ctr = ctr.
			WithSecretVariable("GITLAB_TOKEN", g.PrivateToken).        // for glab
			WithSecretVariable("GITLAB_PRIVATE_TOKEN", g.PrivateToken) // for release-cli
	}
	if g.Repo != "" {
		ctr = ctr.
			WithEnvVariable("GITLAB_REPO", g.Repo).                  // for glab
			WithEnvVariable("CI_PROJECT_ID", url.PathEscape(g.Repo)) // for release-cli
	}
	if g.Group != "" {
		ctr = ctr.WithEnvVariable("GITLAB_GROUP", g.Group) // for glab
	}

	if g.GitDirectory != nil {
		ctr = ctr.WithMountedDirectory("/workspace", g.GitDirectory).
			WithWorkdir("/workspace").
			WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "/workspace"})
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

func hostname(hostOrHostname string) string {
	// if the hostOrHostname is a URL, we want to extract the hostname
	if strings.Contains(hostOrHostname, "://") {
		u, err := url.Parse(hostOrHostname)
		if err != nil {
			return hostOrHostname // return as is if parsing fails
		}
		return u.Hostname()
	}
	return hostOrHostname
}

func serverURL(hostOrHostname string) string {
	// if the hostOrHostname is a URL, we want to extract the server URL
	if strings.Contains(hostOrHostname, "://") {
		u, err := url.Parse(hostOrHostname)
		if err != nil {
			return hostOrHostname // return as is if parsing fails
		}
		return u.Scheme + "://" + u.Host
	}
	return "https://" + hostOrHostname
}
