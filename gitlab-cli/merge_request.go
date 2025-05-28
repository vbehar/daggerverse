package main

import (
	"context"

	"github.com/vbehar/daggerverse/gitlab-cli/internal/dagger"
)

// MergeRequest allows you to interact with GitLab Merge Requests.
func (g *GitlabCli) MergeRequest(
	// ID of the merge request.
	iid string,
) *MergeRequest {
	return &MergeRequest{
		GitlabCli: g,
		Iid:       iid,
	}
}

// MergeRequest allows you to interact with GitLab Merge Requests.
type MergeRequest struct {
	// +private
	GitlabCli *GitlabCli
	Iid       string
}

// Comment adds a comment to the merge request.
// Returns the comment URL.
func (mr *MergeRequest) Comment(
	ctx context.Context,
	// message.
	message string,
	// only add the comment if it doesn't exist yet.
	// +optional
	// +default=false
	unique bool,
	// should we ignore failure?
	// +optional
	// +default=false
	ignoreFailure bool,
) (string, error) {
	if mr.Iid == "" {
		return "", nil
	}

	expect := dagger.ReturnTypeSuccess
	if ignoreFailure {
		expect = dagger.ReturnTypeAny
	}

	var flags []string
	if unique {
		flags = append(flags, "--unique")
	}

	return mr.GitlabCli.Container(ctx).
		WithExec(append([]string{
			"glab", "mr", "note", mr.Iid,
			"--message", message,
		}, flags...), dagger.ContainerWithExecOpts{
			Expect: expect,
		}).
		Stdout(ctx)
}

// Info returns the merge request information.
func (mr *MergeRequest) Info(
	ctx context.Context,
	// output format. Available formats are: text, json.
	// +optional
	// +default="text"
	format string,
	// should we include comments in the output?
	// +optional
	// +default=false
	includeComments bool,
	// should we ignore failure?
	// +optional
	// +default=false
	ignoreFailure bool,
) (string, error) {
	if mr.Iid == "" {
		return "", nil
	}

	expect := dagger.ReturnTypeSuccess
	if ignoreFailure {
		expect = dagger.ReturnTypeAny
	}

	var flags []string
	if includeComments {
		flags = append(flags, "--comments")
	}

	return mr.GitlabCli.Container(ctx).
		WithExec(append([]string{
			"glab", "mr", "view", mr.Iid, "--output", format,
		}, flags...), dagger.ContainerWithExecOpts{
			Expect: expect,
		}).
		Stdout(ctx)
}
