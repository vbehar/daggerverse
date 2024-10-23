// A Dagger Module to extract information about a git reference.
//
// Easily extract information about a git reference (branch, tag, commit hash, committer, commit time, commit message, version),
// and expose it as a JSON file, a directory, or environment variables.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/vbehar/daggerverse/git-info/internal/dagger"
)

// GitInfo contains information about a git reference
type GitInfo struct {
	// git reference used for the git commands
	Ref string
	// branch of the git reference
	Branch string
	// tag of the git reference - if any
	Tag string
	// commit hash of the git reference
	CommitHash string
	// committer information
	CommitUser string
	// commit time
	CommitTime string
	// commit message
	CommitMessage string
	// version of the git reference
	Version string
}

// New returns a new GitInfo instance with information about the git reference
func New(
	ctx context.Context,
	// directory containing the git repository
	// can be either the worktree (including the .git subdirectory)
	// or the .git directory iteself
	gitDirectory *dagger.Directory,
	// git reference to use for the git commands
	// +optional
	// +default="HEAD"
	gitRef string,
	// base container to use for git commands
	// default to cgr.dev/chainguard/git:latest
	// +optional
	gitBaseContainer *dagger.Container,
	// length of the commit hash to use
	// +optional
	// +default=40
	commitHashLength int,
	// format of the commit user to use
	// see https://git-scm.com/docs/git-log#_pretty_formats
	// +optional
	// +default="%an"
	commitUserFormat string,
	// format of the commit time to use
	// see https://git-scm.com/docs/git-log#_pretty_formats
	// +optional
	// +default="%cI"
	commitDateFormat string,
	// format of the commit message to use
	// see https://git-scm.com/docs/git-log#_pretty_formats
	// +optional
	// +default="%B"
	commitMessageFormat string,
) (*GitInfo, error) {
	ctr := gitBaseContainer
	if ctr == nil {
		ctr = dag.Container().From("cgr.dev/chainguard/git:latest")
	}
	ctr = ctr.
		WithMountedDirectory("/workdir", gitDirectory).
		WithWorkdir("/workdir").
		WithExec([]string{"git", "config", "--global", "--add", "safe.directory", "/workdir"})

	branch, err := ctr.
		WithExec([]string{"git", "rev-parse", "--abbrev-ref", gitRef}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch: %w", err)
	}

	tag, _ := ctr.
		WithExec([]string{"git", "describe", "--tags", "--exact-match", gitRef}).
		Stdout(ctx)

	commitHash, err := ctr.
		WithExec([]string{"git", "rev-parse", "--short=" + strconv.Itoa(commitHashLength), gitRef}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit hash with format %q: %w", commitUserFormat, err)
	}

	commitUser, err := ctr.
		WithExec([]string{"git", "show", "-s", "--format=" + commitUserFormat, gitRef}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit user with format %q: %w", commitUserFormat, err)
	}

	commitTime, err := ctr.
		WithExec([]string{"git", "show", "-s", "--format=" + commitDateFormat, gitRef}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit time with format %q: %w", commitDateFormat, err)
	}

	commitMessage, err := ctr.
		WithExec([]string{"git", "show", "-s", "--format=" + commitMessageFormat, gitRef}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit message with format %q: %w", commitMessageFormat, err)
	}

	version, err := ctr.
		WithExec([]string{"git", "describe", "--tags", "--always", gitRef}).
		Stdout(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	gitInfo := &GitInfo{
		Ref:           gitRef,
		Branch:        strings.TrimSpace(branch),
		Tag:           strings.TrimSpace(tag),
		CommitHash:    strings.TrimSpace(commitHash),
		CommitUser:    strings.TrimSpace(commitUser),
		CommitTime:    strings.TrimSpace(commitTime),
		CommitMessage: strings.TrimSpace(commitMessage),
		Version:       strings.TrimSpace(version),
	}
	return gitInfo, nil
}

// Json returns the JSON representation of the git info
func (g *GitInfo) Json() (string, error) { //nolint:stylecheck // we want to name it "json"
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal build infos: %w", err)
	}
	return string(data), nil
}

// JsonFile returns a dagger file containing the JSON representation of the git info
func (g *GitInfo) JsonFile() (*dagger.File, error) { //nolint:stylecheck // we want to name it "json-file"
	data, err := g.Json()
	if err != nil {
		return nil, err
	}
	return dag.Directory().
		WithNewFile("git-info.json", data).
		File("git-info.json"), nil
}

// Directory returns a dagger directory containing the git info - each field is stored in a separate file:
// - ref: git reference used for the git commands
// - branch: branch of the git reference
// - tag: tag of the git reference - if any
// - commit-hash: commit hash of the git reference
// - commit-user: committer information
// - commit-time: commit time
// - commit-message: commit message
// - version: version of the git reference
func (g *GitInfo) Directory() *dagger.Directory {
	return dag.Directory().
		WithNewFile("ref", g.Ref).
		WithNewFile("branch", g.Branch).
		WithNewFile("tag", g.Tag).
		WithNewFile("commit-hash", g.CommitHash).
		WithNewFile("commit-user", g.CommitUser).
		WithNewFile("commit-time", g.CommitTime).
		WithNewFile("commit-message", g.CommitMessage).
		WithNewFile("version", g.Version)
}

// SetEnvVariablesOnContainer sets the git info as environment variables on the container:
// - GIT_REF: git reference used for the git commands
// - GIT_BRANCH: branch of the git reference
// - GIT_TAG: tag of the git reference - if any
// - GIT_COMMIT_HASH: commit hash of the git reference
// - GIT_COMMIT_USER: committer information
// - GIT_COMMIT_TIME: commit time
// - GIT_COMMIT_MESSAGE: commit message
// - GIT_VERSION: version of the git reference
func (g *GitInfo) SetEnvVariablesOnContainer(
	ctr *dagger.Container,
) *dagger.Container {
	return ctr.
		WithEnvVariable("GIT_REF", g.Ref).
		WithEnvVariable("GIT_BRANCH", g.Branch).
		WithEnvVariable("GIT_TAG", g.Tag).
		WithEnvVariable("GIT_COMMIT_HASH", g.CommitHash).
		WithEnvVariable("GIT_COMMIT_USER", g.CommitUser).
		WithEnvVariable("GIT_COMMIT_TIME", g.CommitTime).
		WithEnvVariable("GIT_COMMIT_MESSAGE", g.CommitMessage).
		WithEnvVariable("GIT_VERSION", g.Version)
}
