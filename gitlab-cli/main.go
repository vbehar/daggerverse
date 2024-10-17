// A Dagger Module to install and run the GitLab CLI.
//
// GitLab CLI is a command-line tool that allows you to interact with the GitLab API.
// See https://gitlab.com/gitlab-org/cli for more information.
package main

import (
	"context"

	"github.com/vbehar/daggerverse/gitlab-cli/internal/dagger"
)

// GitlabCli is a Dagger Module to interact with the GitLab CLI.
type GitlabCli struct {
	Token *dagger.Secret
	Host  string
	Repo  string
	Group string
}

func New(
	// token to use for authentication with GitLab.
	// +optional
	token *dagger.Secret,
	// host of the GitLab instance.
	// +optional
	host string,
	// default gitlab repository for commands accepting the --repo flag.
	// +optional
	repo string,
	// default gitlab group for commands accepting the --group flag.
	// +optional
	group string,
) *GitlabCli {
	return &GitlabCli{
		Token: token,
		Host:  host,
		Repo:  repo,
		Group: group,
	}
}

// Container returns a container with the GitLab CLI installed.
func (g *GitlabCli) Container() *dagger.Container {
	ctr := dag.Container().
		From("cgr.dev/chainguard/wolfi-base").
		WithExec([]string{"apk", "add", "--update", "--no-cache",
			"ca-certificates",
			"glab",
		})

	if g.Token != nil {
		ctr = ctr.WithSecretVariable("GITLAB_TOKEN", g.Token)
	}
	if g.Host != "" {
		ctr = ctr.WithEnvVariable("GITLAB_HOST", g.Host)
	}
	if g.Repo != "" {
		ctr = ctr.WithEnvVariable("GITLAB_REPO", g.Repo)
	}
	if g.Group != "" {
		ctr = ctr.WithEnvVariable("GITLAB_GROUP", g.Group)
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
		ctr = g.Container()
	}

	return ctr.
		WithEntrypoint([]string{"glab"}).
		WithExec(args, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		Stdout(ctx)
}
