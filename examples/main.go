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
	version Optional[string],
	logLevel Optional[string],
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
		dag.Host().Directory("testdata"),
		version.GetOr("v0.0.1"),
		repoName,
		ArtifactoryPublishGoLibOpts{
			LogLevel: logLevel.GetOr("debug"),
		},
	)
}
