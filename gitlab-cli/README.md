# GitLab CLI Dagger Module

This is a [Dagger](https://dagger.io/) module for the [GitLab CLI](https://gitlab.com/gitlab-org/cli) tool.

Use it to install the glab CLI in a Dagger container, and easily interact with GitLab API.

Read the documentation at <https://daggerverse.dev/mod/github.com/vbehar/daggerverse/gitlab-cli>.

## Usage

Run a basic glab command:

```bash
$ dagger call -m github.com/vbehar/daggerverse/gitlab-cli \
	glab --args version
```

Or a more complex command, to list the releases of a specific repository in a self-hosted GitLab instance:

```bash
$ dagger call -m github.com/vbehar/daggerverse/gitlab-cli \
	--token=env:GITLAB_TOKEN \
	--host=https://gitlab.example.com \
	--repo=my-owner/my-project \
	glab \
	--args "release,list"
```
