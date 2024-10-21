# GitLab CLI Dagger Module

This is a [Dagger](https://dagger.io/) module for the [GitLab CLI](https://gitlab.com/gitlab-org/cli) tool.
It also contains the (deprecated) [GitLab Release CLI](https://gitlab.com/gitlab-org/release-cli) tool.

Use it to install the glab CLI in a Dagger container, and easily interact with GitLab API.

Read the documentation at <https://daggerverse.dev/mod/github.com/vbehar/daggerverse/gitlab-cli>.

## Usage

Run a basic glab command:

```bash
$ dagger call -m github.com/vbehar/daggerverse/gitlab-cli \
	run --args version
```

Or a more complex command, to list the releases of a specific repository in a self-hosted GitLab instance:

```bash
$ dagger call -m github.com/vbehar/daggerverse/gitlab-cli \
	--private-token=env:GITLAB_TOKEN \
	--host=https://gitlab.example.com \
	--repo=my-owner/my-project \
	run \
	--args "release,list"
```

### From inside a GitLab CI pipeline

Note that the [GitLab CLI](https://gitlab.com/gitlab-org/cli) tool [doesn't support "job tokens"](https://gitlab.com/gitlab-org/cli/-/issues/1220) (i.e. the `CI_JOB_TOKEN` environment variable) yet.
Instead, you'll need to use the [GitLab Release CLI](https://gitlab.com/gitlab-org/release-cli) tool - which properly supports job tokens, but only to manipulate releases.

For example, to create a new release (and tag) in a GitLab CI environment, you can use the following command:

```bash
$ dagger call -m github.com/vbehar/daggerverse/gitlab-cli \
	--job-token=env:CI_JOB_TOKEN \
	--host=https://gitlab.example.com \
	--repo=my-owner/my-project \
	release \
		--tag-name=v1.2.3 \
		--description-file=changelog.md \
		create \
			--git-ref=main
```
