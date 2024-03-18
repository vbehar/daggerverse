# JFrog Artifactory Dagger Module

This is a [Dagger](https://dagger.io/) module to interact with [JFrog Artifactory](https://jfrog.com/artifactory/).

## Usage

Ping your artifactory instance:

```bash
$ dagger call --instance-url=https://artifactory.example.com/artifactory command --cmd rt,ping stdout
┃ OK
```

Publish a Go library:

```bash
$ export ARTIFACTORY_PASSWORD=xyz
$ dagger call --instance-url=https://artifactory.example.com/artifactory --username=YOUR_USER --password=${ARTIFACTORY_PASSWORD} publishGoLib --repo YOUR_ARTIFACTORY_REPO --src ./testdata --version v0.0.1 --log-level debug
```
