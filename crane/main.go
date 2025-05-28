// A Dagger Module to install and run the Crane CLI.
//
// Crane is a command-line tool that allows you to interact with container registries.
// See https://github.com/google/go-containerregistry/tree/main/cmd/crane for more information.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/vbehar/daggerverse/crane/internal/dagger"
)

const (
	// use fixed base images for reproductible builds and improved caching
	// the base crane image: https://images.chainguard.dev/directory/image/crane/overview
	// retrieve the latest sha256 hash with: `crane digest cgr.dev/chainguard/crane:latest`
	// and to retrieve its creation time: `crane config cgr.dev/chainguard/crane:latest | jq .created`
	// This one is from 2025-05-22T20:15:28Z
	baseCraneImage = "cgr.dev/chainguard/crane:latest@sha256:bd1a13a9265e46025f6d1393f7b75f5678ad74381ea56f110b04540045d97db3"
)

// Crane is a Dagger Module to interact with the Crane CLI.
type Crane struct {
	Registry string
	Username string
	Password *dagger.Secret
	Insecure bool
	Platform string
}

func New(
	// registry to authenticate to
	// +optional
	registry string,
	// username to use for authentication with the registry
	// +optional
	username string,
	// password to use for authentication with the registry
	// +optional
	password *dagger.Secret,
	// allow insecure connections to the registry
	// +optional
	insecure bool,
	// platform to request when listing images
	// default to all platforms
	// +optional
	platform string,
) *Crane {
	return &Crane{
		Registry: registry,
		Username: username,
		Password: password,
		Insecure: insecure,
		Platform: platform,
	}
}

// Login returns a new Crane instance with the given registry and credentials.
func (c *Crane) Login(
	// registry to authenticate to
	registry string,
	// username to use for authentication with the registry
	username string,
	// password to use for authentication with the registry
	password *dagger.Secret,
) *Crane {
	return &Crane{
		Registry: registry,
		Username: username,
		Password: password,
		Insecure: c.Insecure,
		Platform: c.Platform,
	}
}

// WithPlatform returns a new Crane instance with the given platform.
// This is useful when you want to list images for a specific platform.
// If the platform is empty, it defaults to all platforms.
func (c *Crane) WithPlatform(
	// platform to request when listing images
	platform string,
) *Crane {
	return &Crane{
		Registry: c.Registry,
		Username: c.Username,
		Password: c.Password,
		Insecure: c.Insecure,
		Platform: platform,
	}
}

// Container returns a container with the Crane CLI installed
// and the registry configured - if a registry and credentials are provided.
func (c *Crane) Container() *dagger.Container {
	ctr := dag.Container().From(baseCraneImage)

	if c.Registry != "" {
		ctr = ctr.WithEnvVariable("REGISTRY_HOST", c.Registry)
	}
	if c.Username != "" && c.Password != nil {
		ctr = ctr.
			WithEnvVariable("REGISTRY_USERNAME", c.Username).
			WithSecretVariable("REGISTRY_PASSWORD", c.Password).
			WithExec([]string{
				"/bin/sh", "-c",
				"crane auth login --username $REGISTRY_USERNAME --password $REGISTRY_PASSWORD $REGISTRY_HOST",
			})
	}
	return ctr
}

// Run runs the crane CLI with the given arguments.
func (c *Crane) Run(
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
		ctr = c.Container()
	}

	if c.Platform != "" {
		args = append([]string{"--platform", c.Platform}, args...)
	}
	if c.Insecure {
		args = append([]string{"--insecure"}, args...)
	}

	return ctr.
		WithEntrypoint([]string{"crane"}).
		WithExec(args, dagger.ContainerWithExecOpts{
			UseEntrypoint: true,
		}).
		Stdout(ctx)
}

// Ls lists the images in the given repository.
func (c *Crane) Ls(
	ctx context.Context,
	// repository to list images from
	repository string,
	// print the full image reference
	// +optional
	// +default=false
	fullRef bool,
	// omit digest tags (e.g., ':sha256-...')
	// +optional
	// +default=false
	omitDigestTags bool,
	// +optional
	ctr *dagger.Container,
) ([]string, error) {
	args := []string{
		"ls",
		repository,
	}
	if fullRef {
		args = append(args, "--full-ref")
	}
	if omitDigestTags {
		args = append(args, "--omit-digest-tags")
	}

	output, err := c.Run(ctx, args, ctr)
	if err != nil {
		return nil, fmt.Errorf("failed to run crane ls: %w", err)
	}

	result := strings.Split(output, "\n")
	return result, nil
}

// ImageTagExists checks if the given image tag exists.
func (c *Crane) ImageTagExists(
	ctx context.Context,
	// image to check
	// format: <repository>:<tag>
	image string,
	// +optional
	ctr *dagger.Container,
) (bool, error) {
	repository, tag, ok := strings.Cut(image, ":")
	if !ok {
		return false, fmt.Errorf("invalid image format: %s", image)
	}

	allTags, err := c.Ls(ctx, repository, false, true, ctr)
	if err != nil {
		return false, fmt.Errorf("failed to list tags: %w", err)
	}

	for _, t := range allTags {
		if t == tag {
			return true, nil
		}
	}
	return false, nil
}
