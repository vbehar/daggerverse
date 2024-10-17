# Git Changelog Dagger Module

This is a [Dagger](https://dagger.io/) module for the [git-chglog](https://github.com/git-chglog/git-chglog) tool.

Use it to install the git-chglog CLI in a Dagger container, and easily generate changelogs from your Git repositories.

Read the documentation at <https://daggerverse.dev/mod/github.com/vbehar/daggerverse/git-chglog>.

## Usage

Generate a changelog with the default config/template:

```bash
$ dagger call -m github.com/vbehar/daggerverse/git-chglog \
	changelog \
	  --git-repository=.git \
	  --version=v1.2.3 \
	  contents
```

Or with a custom config file and template:

```bash
$ dagger call -m github.com/vbehar/daggerverse/git-chglog \
	--chglog-dir=.chglog \
	changelog \
	  --git-repository=.git \
	  --version=v1.2.3 \
	  contents
```
