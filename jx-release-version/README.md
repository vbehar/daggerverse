# JX Release Version Dagger Module

This is a [Dagger](https://dagger.io/) module for the [JX Release Version](https://github.com/jenkins-x-plugins/jx-release-version) tool.

Use it to install the jx-release-version CLI in a Dagger container, and easily calculate the next release version of your Git repository.

Read the documentation at <https://daggerverse.dev/mod/github.com/vbehar/daggerverse/jx-release-version>.

## Usage

Get the next release version of your Git repository:

```bash
$ dagger call -m github.com/vbehar/daggerverse/jx-release-version \
	next-version --git-directory=/path/to/.git
```

Automatically create and push a new Git tag with the next release version:

```bash
$ dagger call -m github.com/vbehar/daggerverse/jx-release-version \
	tag \
	--git-directory=/path/to/.git \
	--git-token=env:GIT_TOKEN
```
