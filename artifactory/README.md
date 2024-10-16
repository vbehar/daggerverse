# JFrog Artifactory Dagger Module

This is a [Dagger](https://dagger.io/) module to interact with [JFrog Artifactory](https://jfrog.com/artifactory/).

Use it to upload artifacts to Artifactory.

Read the documentation at <https://daggerverse.dev/mod/github.com/vbehar/daggerverse/artifactory>.

## Usage

Ping your artifactory instance:

```bash
$ dagger call -i -m github.com/vbehar/daggerverse/artifactory \
    --instance-url=https://artifactory.example.com/artifactory \
    command --cmd rt,ping \
    stdout
┃ OK
```

Publish a Go library:

```bash
$ export ARTIFACTORY_USER=YOUR_USER
$ export ARTIFACTORY_PASSWORD=xyz
$ export ARTIFACTORY_REPO=YOUR_ARTIFACTORY_REPO
$ dagger call -i -m github.com/vbehar/daggerverse/artifactory \
    --instance-url=https://artifactory.example.com/artifactory --username=${ARTIFACTORY_USER} --password=env:ARTIFACTORY_PASSWORD \
    publish-go-lib --repo ${ARTIFACTORY_REPO} --src ./testdata --version v0.0.1 \
    stdout
```
