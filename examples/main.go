package main

import (
	"context"
)

type Examples struct{}

func (e *Examples) PublishGoLibToArtifactory(
	ctx context.Context,
	instanceName string,
	artifactoryUser string,
	artifactoryPassword *Secret,
	// +optional
	// +default="v0.0.1"
	version string,
	// +optional
	// +default="debug"
	logLevel string,
) *Container {
	var (
		instanceURL = "https://artifactory." + instanceName + ".org/artifactory"
		repoName    = "go-snapshot-" + instanceName
	)

	return dag.Artifactory(instanceURL, ArtifactoryOpts{
		InstanceName: instanceName,
		Username:     artifactoryUser,
		Password:     artifactoryPassword,
	}).PublishGoLib(
		dag.CurrentModule().Source().Directory("testdata"),
		version,
		repoName,
		ArtifactoryPublishGoLibOpts{
			LogLevel: logLevel,
		},
	)
}
