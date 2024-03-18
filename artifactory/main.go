package main

import "strings"

const (
	defaultInstanceName = "default"
	defaultGenericImage = "cgr.dev/chainguard/wolfi-base"
	defaultGoImage      = "cgr.dev/chainguard/go:latest-dev"
)

type Artifactory struct {
	InstanceName    string
	InstanceURL     string
	Username        string
	Password        *Secret
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
	password *Secret,
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
	// container to configure.
	ctr *Container,
) *Container {
	ctr = dag.Jfrogcli(JfrogcliOpts{
		Version: a.JfrogCliVersion,
	}).Install(JfrogcliInstallOpts{
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
			}, ContainerWithExecOpts{
				SkipEntrypoint: true,
			})
	}

	return ctr.
		WithEnvVariable("ARTIFACTORY_URL", a.InstanceURL).
		WithEnvVariable("ARTIFACTORY_USERNAME", a.Username).
		WithSecretVariable("ARTIFACTORY_PASSWORD", a.Password).
		WithExec([]string{
			"/bin/sh", "-c",
			"echo ${ARTIFACTORY_PASSWORD} | jf config add --artifactory-url ${ARTIFACTORY_URL} --user ${ARTIFACTORY_USERNAME} --password-stdin --overwrite " + a.InstanceName,
		}, ContainerWithExecOpts{
			SkipEntrypoint: true,
		}).
		WithoutEnvVariable("ARTIFACTORY_URL").
		WithoutEnvVariable("ARTIFACTORY_USERNAME").
		WithoutEnvVariable("ARTIFACTORY_PASSWORD")
}

func configureArtifactory(a *Artifactory) WithContainerFunc {
	return func(ctr *Container) *Container {
		return a.Configure(ctr)
	}
}

// Command runs the given artifactory (jf) command in the given container.
func (a *Artifactory) Command(
	// jf command to run. the "jf" prefix will be added automatically.
	cmd []string,
	// container to run the command in. If empty, a new container will be created.
	// +optional
	ctr *Container,
	// log level to use for the command. If empty, the default log level will be used.
	// +optional
	logLevel string,
) *Container {
	if ctr == nil {
		ctr = dag.Container().From(defaultGenericImage)
	}
	return ctr.
		With(configureArtifactory(a)).
		With(jfLogLevel(logLevel)).
		WithFocus().
		WithExec(append([]string{"jf"}, cmd...), ContainerWithExecOpts{SkipEntrypoint: true}).
		WithoutFocus()
}

func jfLogLevel(
	// +optional
	logLevel string,
) WithContainerFunc {
	return func(ctr *Container) *Container {
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
) WithContainerFunc {
	return func(ctr *Container) *Container {
		return a.Command(cmd, ctr, logLevel)
	}
}

// PublishGoLib publishes a Go library to the given repository.
func (a *Artifactory) PublishGoLib(
	// directory containing the Go library to publish.
	src *Directory,
	// version of the library to publish.
	version string,
	// name of the repository to publish to.
	repo string,
	// log level to use for the command. If empty, the default log level will be used.
	// +optional
	logLevel string,
) *Container {
	return dag.Container().From(defaultGoImage).
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
