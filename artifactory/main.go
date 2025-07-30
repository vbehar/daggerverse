// A Dagger module to interact with JFrog Artifactory.
//
// Artifactory is a service that provides repositories for storing and managing artifacts.
package main

import (
	"context"
	"strings"

	"github.com/vbehar/daggerverse/artifactory/internal/dagger"
)

const (
	defaultInstanceName = "default"

	// use fixed base images for reproductible builds and improved caching
	// the base wolfi image: https://images.chainguard.dev/directory/image/wolfi-base/overview
	// retrieve the latest sha256 hash with: `crane digest cgr.dev/chainguard/wolfi-base:latest`
	// and to retrieve its creation time: `crane config cgr.dev/chainguard/wolfi-base:latest | jq .created`
	// This one is from 2025-06-02T17:31:02Z
	baseWolfiImage = "cgr.dev/chainguard/wolfi-base:latest@sha256:57428116d2d7c27d1d4de4103e19b40bb8d2942ff6dff31b900e55efedeb7e30"

	// the base BUILD image: https://images.chainguard.dev/directory/image/go/overview
	// retrieve the latest sha256 hash with: `crane digest cgr.dev/chainguard/go:latest-dev`
	// and to retrieve its creation time: `crane config cgr.dev/chainguard/go:latest-dev | jq .created`
	// This one is from 2025-06-06T21:28:17Z
	baseGoImage = "cgr.dev/chainguard/go:latest-dev@sha256:faa589370de5c382cb7c4ae7313bd0fa677db4b70ae72013307d7fc93890e272"
)

// Artifactory is a Dagger Module to interact with JFrog Artifactory.
type Artifactory struct {
	// name of the Artifactory instance.
	InstanceName string
	// URL of the Artifactory instance.
	InstanceURL string
	// username to use for authentication. If empty, authentication will not be configured.
	Username string
	// password to use for authentication.
	Password *dagger.Secret
	// version of the JFrog CLI.
	JfrogCliVersion string
}

func New(
	// URL of the Artifactory instance.
	instanceURL string,
	// username to use for authentication. If empty, authentication will not be configured.
	// +optional
	username string,
	// password to use for authentication.
	// +optional
	password *dagger.Secret,
	// name of the Artifactory instance to configure. Defaults to "default".
	// +optional
	// +default="default"
	instanceName string,
	// version of the JFrog CLI to install. If empty, the latest version will be installed.
	// +optional
	jfrogCliVersion string,
) *Artifactory {
	return &Artifactory{
		InstanceName:    instanceName,
		InstanceURL:     instanceURL,
		Username:        username,
		Password:        password,
		JfrogCliVersion: jfrogCliVersion,
	}
}

// Configure configures the given container to use the Artifactory instance.
func (a *Artifactory) Configure(
	// container to configure. If empty, a new container will be created.
	// +optional
	ctr *dagger.Container,
) *dagger.Container {
	ctr = dag.Jfrogcli(dagger.JfrogcliOpts{
		Version: a.JfrogCliVersion,
	}).Install(dagger.JfrogcliInstallOpts{
		Base: ctr,
	})

	if a.Username == "" || a.Password == nil {
		return ctr.
			WithExec([]string{
				"jf",
				"config", "add",
				"--artifactory-url", a.InstanceURL,
				"--overwrite",
				a.InstanceName,
			})
	}

	return ctr.
		WithEnvVariable("ARTIFACTORY_URL", a.InstanceURL).
		WithEnvVariable("ARTIFACTORY_USERNAME", a.Username).
		WithSecretVariable("ARTIFACTORY_PASSWORD", a.Password).
		WithExec([]string{
			"/bin/sh", "-c",
			"echo ${ARTIFACTORY_PASSWORD} | jf config add --artifactory-url ${ARTIFACTORY_URL} --user ${ARTIFACTORY_USERNAME} --password-stdin --overwrite " + a.InstanceName,
		}).
		WithoutEnvVariable("ARTIFACTORY_URL").
		WithoutEnvVariable("ARTIFACTORY_USERNAME").
		WithoutSecretVariable("ARTIFACTORY_PASSWORD")
}

// Command runs the given artifactory (jf) command in the given container.
func (a *Artifactory) Command(
	// jf command to run. the "jf" prefix will be added automatically.
	cmd []string,
	// container to run the command in. If empty, a new container will be created.
	// +optional
	ctr *dagger.Container,
	// log level to use for the command. If empty, the default log level will be used.
	// +optional
	logLevel string,
) *dagger.Container {
	if ctr == nil {
		ctr = dag.Container().From(baseWolfiImage)
	}
	return ctr.
		With(configureArtifactory(a)).
		With(jfLogLevel(logLevel)).
		WithExec(append([]string{"jf"}, cmd...))
}

// PublishFile publishes a single file to artifactory.
func (a *Artifactory) PublishFile(
	ctx context.Context,
	// file to publish.
	file *dagger.File,
	// target path in artifactory.
	destination string,
	// log level to use for the command. If empty, the default log level will be used.
	// +optional
	logLevel string,
) (string, error) {
	return a.Command(
		[]string{
			"rt", "u",
			"/src",
			destination,
		},
		dag.Container().From(baseWolfiImage).
			WithMountedFile("/src", file),
		logLevel).
		Stdout(ctx)
}

// PublishGoLib publishes a Go library to the given repository.
func (a *Artifactory) PublishGoLib(
	ctx context.Context,
	// directory containing the Go library to publish.
	src *dagger.Directory,
	// version of the library to publish.
	// Default to the "git" version (from the `git describe` cmd).
	// +optional
	version string,
	// name of the repository to publish to.
	repo string,
	// log level to use for the command. If empty, the default log level will be used.
	// +optional
	logLevel string,
) *dagger.Container {
	if version == "" {
		var err error
		version, err = dag.GitInfo(src).Version(ctx)
		if err != nil {
			version = "v0.0.1"
		}
		version = strings.TrimSpace(version)
	}
	return dag.Container().From(baseGoImage).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithEnvVariable("GOWORK", "off"). // jf tries to run `go list -mod=mod -m` which won't work in workspace mode
		With(jfCommand(a, []string{
			"go-config",
			"--repo-deploy=" + repo,
			"--server-id-deploy=" + a.InstanceName,
		}, "")).
		With(jfCommand(a, []string{
			"go-publish",
			"--detailed-summary",
			version,
		}, logLevel)).
		WithoutEnvVariable("GOWORK")
}

func configureArtifactory(a *Artifactory) dagger.WithContainerFunc {
	return func(ctr *dagger.Container) *dagger.Container {
		return a.Configure(ctr)
	}
}

func jfLogLevel(
	// +optional
	logLevel string,
) dagger.WithContainerFunc {
	return func(ctr *dagger.Container) *dagger.Container {
		if logLevel != "" {
			ctr = ctr.WithEnvVariable("JFROG_CLI_LOG_LEVEL", strings.ToUpper(logLevel))
		}
		return ctr
	}
}

func jfCommand(
	a *Artifactory,
	cmd []string,
	// +optional
	logLevel string,
) dagger.WithContainerFunc {
	return func(ctr *dagger.Container) *dagger.Container {
		return a.Command(cmd, ctr, logLevel)
	}
}
