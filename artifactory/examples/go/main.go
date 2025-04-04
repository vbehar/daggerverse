package main

import (
	"context"

	"github.com/vbehar/daggerverse/artifactory/examples/go/internal/dagger"
)

type Examples struct{}

func (e *Examples) Artifactory_PublishFile(
	ctx context.Context,
	instanceName string,
	artifactoryUser string,
	artifactoryPassword *dagger.Secret,
	// +optional
	// +default="debug"
	logLevel string,
) (string, error) {
	instanceURL := "https://artifactory." + instanceName + ".org/artifactory"

	return dag.Artifactory(instanceURL, dagger.ArtifactoryOpts{
		InstanceName: instanceName,
		Username:     artifactoryUser,
		Password:     artifactoryPassword,
	}).PublishFile(
		ctx,
		dag.CurrentModule().Source().Directory("testdata").File("main.go"),
		"some-repo/some/path/main.go",
		dagger.ArtifactoryPublishFileOpts{
			LogLevel: logLevel,
		},
	)
}

func (e *Examples) Artifactory_PublishGoLib(
	ctx context.Context,
	instanceName string,
	artifactoryUser string,
	artifactoryPassword *dagger.Secret,
	// +optional
	// +default="v0.0.1"
	version string,
	// +optional
	// +default="debug"
	logLevel string,
) *dagger.Container {
	var (
		instanceURL = "https://artifactory." + instanceName + ".org/artifactory"
		repoName    = "go-snapshot-" + instanceName
	)

	return dag.Artifactory(instanceURL, dagger.ArtifactoryOpts{
		InstanceName: instanceName,
		Username:     artifactoryUser,
		Password:     artifactoryPassword,
	}).PublishGoLib(
		dag.CurrentModule().Source().Directory("testdata"),
		repoName,
		dagger.ArtifactoryPublishGoLibOpts{
			Version:  version,
			LogLevel: logLevel,
		},
	)
}
