package main

import (
	"context"
	"fmt"

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

// Diff returns the merge request diff.
func (mr *MergeRequest) Diff(
	ctx context.Context,
	// use raw diff format?
	// +optional
	// +default=false
	raw bool,
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
	if raw {
		flags = append(flags, "--raw")
	}

	return mr.GitlabCli.Container(ctx).
		WithExec(append([]string{
			"glab", "mr", "diff", mr.Iid,
		}, flags...), dagger.ContainerWithExecOpts{
			Expect: expect,
		}).
		Stdout(ctx)
}

// Commits returns the merge request commits.
// In JSON format
// https://docs.gitlab.com/api/merge_requests/#get-single-merge-request-commits
func (mr *MergeRequest) Commits(
	ctx context.Context,
	// should we ignore failure?
	// +optional
	// +default=false
	ignoreFailure bool,
) (string, error) {
	if mr.Iid == "" {
		return "", nil
	}

	endpoint := fmt.Sprintf("/projects/:fullpath/merge_requests/%s/commits", mr.Iid)

	expect := dagger.ReturnTypeSuccess
	if ignoreFailure {
		expect = dagger.ReturnTypeAny
	}

	return mr.GitlabCli.Container(ctx).
		WithExec([]string{
			"glab", "api", endpoint,
		}, dagger.ContainerWithExecOpts{
			Expect: expect,
		}).
		Stdout(ctx)
}

// Update updates a merge request.
func (mr *MergeRequest) Update(
	ctx context.Context,
	// the title. If empty, the title will not be updated.
	// +optional
	title string,
	// the description. If empty, the description will not be updated.
	// +optional
	description string,
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
	if title != "" {
		flags = append(flags, "--title", title)
	}
	if description != "" {
		flags = append(flags, "--description", description)
	}

	return mr.GitlabCli.Container(ctx).
		WithExec(append([]string{
			"glab", "mr", "update", mr.Iid,
		}, flags...), dagger.ContainerWithExecOpts{
			Expect: expect,
		}).
		Stdout(ctx)
}
