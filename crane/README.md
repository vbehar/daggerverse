# Crane CLI Dagger Module

This is a [Dagger](https://dagger.io/) module for the [Crane](https://github.com/google/go-containerregistry/tree/main/cmd/crane) tool.

Use it to install the crane CLI in a Dagger container, and easily interact with container registries.

Read the documentation at <https://daggerverse.dev/mod/github.com/vbehar/daggerverse/crane>.

## Usage

Run a basic crane command, to list the tags of a public repository:

```bash
$ dagger call -m github.com/vbehar/daggerverse/crane \
	run --args ls,registry.dagger.io/engine
```

Or a more complex command, to list the tags of a private repository:

```bash
$ dagger call -m github.com/vbehar/daggerverse/crane \
	--registry=docker.artifactory.example.com \
	--username=vbehar \
	--token=env:ARTIFACTORY_API_KEY \
	run \
	--args ls,docker.artifactory.example.com/my-image
```

You can also easily check if an image tag exists:

```bash
$ dagger call -m github.com/vbehar/daggerverse/crane \
	image-tag-exists --image=registry.dagger.io/engine:v0.13.5
```
