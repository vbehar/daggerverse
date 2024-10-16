// A Dagger Module to interact with the jx-release-version CLI.
//
// jx-release-version is a CLI tool that helps you to release new versions of your projects.
// See https://github.com/jenkins-x-plugins/jx-release-version for more information.
package main

import (
	"context"
	"regexp"
	"strconv"

	"github.com/vbehar/daggerverse/jx-release-version/internal/dagger"
)

// JxReleaseVersion is a Dagger Module to interact with the jx-release-version CLI.
type JxReleaseVersion struct {
	ImageRepository string
	ImageTag        string
	LogLevel        string
}

func New(
	// image repository of jx-release-version.
	// +optional
	// +default="ghcr.io/jenkins-x/jx-release-version"
	imageRepository string,
	// image tag of jx-release-version.
	// See https://github.com/jenkins-x-plugins/jx-release-version/releases for available tags.
	// +optional
	// +default="2.7.6@sha256:603d9c7c3cdbb14210abbd04d386f138c0490ebce267f341053c5b24f32fa772"
	imageTag string,
	// log level to use for the command.
	// +optional
	// +default="info"
	logLevel string,
) *JxReleaseVersion {
	return &JxReleaseVersion{
		ImageRepository: imageRepository,
		ImageTag:        imageTag,
		LogLevel:        logLevel,
	}
}

// Container returns a container with jx-release-version installed.
func (jx *JxReleaseVersion) Container(
	ctx context.Context,
	// git directory to include in the container.
	// +optional
	gitDirectory *dagger.Directory,
) *dagger.Container {
	ctr := dag.Container().
		From(jx.ImageRepository+":"+jx.ImageTag).
		WithEnvVariable("JX_LOG_LEVEL", jx.LogLevel)

	if gitDirectory != nil {
		ctr = ctr.
			WithMountedDirectory("/workspace",
				gitDirectory.With(httpsInsteadOfGit(ctx)),
			).
			WithWorkdir("/workspace")
	}

	return ctr
}

// NextVersion returns the next version of the given git repository.
func (jx *JxReleaseVersion) NextVersion(
	ctx context.Context,
	// git directory to include in the container.
	gitDirectory *dagger.Directory,
	// If true, fetch tags from the remote repository before detecting the previous version.
	// +optional
	// +default=false
	fetchTags bool,
	// strategy to use to read the previous version.
	// See https://github.com/jenkins-x-plugins/jx-release-version for doc.
	// +optional
	// +default="auto"
	previousVersionStrategy string,
	// strategy to use to calculate the next version.
	// See https://github.com/jenkins-x-plugins/jx-release-version for doc.
	// +optional
	// +default="auto"
	nextVersionStrategy string,
	// output format of the next version.
	// See https://github.com/jenkins-x-plugins/jx-release-version for doc.
	// +optional
	// +default="{{.Major}}.{{.Minor}}.{{.Patch}}"
	outputFormat string,
) (string, error) {
	return jx.Container(ctx, gitDirectory).
		WithFocus().
		WithExec([]string{
			"jx-release-version",
			"--fetch-tags=" + strconv.FormatBool(fetchTags),
			"--previous-version=" + previousVersionStrategy,
			"--next-version=" + nextVersionStrategy,
			"--output-format=" + outputFormat,
		}).
		Stdout(ctx)
}

// Tag tags the current commit with the next version - and pushes the tag to the remote repository.
// It returns the next version.
func (jx *JxReleaseVersion) Tag(
	ctx context.Context,
	// git directory to include in the container.
	gitDirectory *dagger.Directory,
	// If true, fetch tags from the remote repository before detecting the previous version.
	// +optional
	// +default=false
	fetchTags bool,
	// strategy to use to read the previous version.
	// See https://github.com/jenkins-x-plugins/jx-release-version for doc.
	// +optional
	// +default="auto"
	previousVersionStrategy string,
	// strategy to use to calculate the next version.
	// See https://github.com/jenkins-x-plugins/jx-release-version for doc.
	// +optional
	// +default="auto"
	nextVersionStrategy string,
	// output format of the next version.
	// See https://github.com/jenkins-x-plugins/jx-release-version for doc.
	// +optional
	// +default="{{.Major}}.{{.Minor}}.{{.Patch}}"
	outputFormat string,
	// prefix for the new tag - prefixed before the output.
	// +optional
	// +default="v"
	tagPrefix string,
	// git token to use for authentication when pushing the tag to the remote.
	gitToken *dagger.Secret,
	// name of the author/committer used to create the tag.
	// +optional
	// +default="jx-release-version"
	gitUser string,
	// email of the author/committer used to create the tag.
	// +optional
	// +default="jx-release-version@jenkins-x.io"
	gitEmail string,
	// if true, push the tag to the remote.
	// +optional
	// +default=true
	pushTag bool,
) (string, error) {
	return jx.Container(ctx, gitDirectory).
		WithSecretVariable("GIT_TOKEN", gitToken).
		WithFocus().
		WithExec([]string{
			"jx-release-version",
			"--fetch-tags=" + strconv.FormatBool(fetchTags),
			"--previous-version=" + previousVersionStrategy,
			"--next-version=" + nextVersionStrategy,
			"--output-format=" + outputFormat,
			"--tag-prefix=" + tagPrefix,
			"--git-user=" + gitUser,
			"--git-email=" + gitEmail,
			"--push-tag=" + strconv.FormatBool(pushTag),
			"--tag",
		}).
		Stdout(ctx)
}

// httpsInsteadOfGit modifies the git config to use https instead of git.
func httpsInsteadOfGit(ctx context.Context) dagger.WithDirectoryFunc {
	const configFilePath = "config"
	gitHostnameRegexp := regexp.MustCompile(`git@([^:]+):`)

	return func(gitDir *dagger.Directory) *dagger.Directory {
		originalGitConfig, err := gitDir.File(configFilePath).Contents(ctx)
		if err != nil {
			return gitDir
		}

		gitConfig := gitHostnameRegexp.ReplaceAllStringFunc(originalGitConfig, func(match string) string {
			hostname := gitHostnameRegexp.FindStringSubmatch(match)[1]
			return "https://" + hostname + "/"
		})
		if gitConfig == originalGitConfig {
			return gitDir
		}

		return gitDir.
			WithoutFile(configFilePath).
			WithNewFile(configFilePath, gitConfig)
	}
}
