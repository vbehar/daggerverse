# JFrog Artifactory Dagger Module

This is a [Dagger](https://dagger.io/) module to interact with [JFrog Artifactory](https://jfrog.com/artifactory/).

## Usage

Ping your artifactory instance:

```bash
$ dagger call --instance-url=https://artifactory.example.com/artifactory command --cmd rt,ping --log-level=debug
✔ dagger call command [5.95s]
┃ OK
✔ exec jf rt ping [0.07s]
┃ 09:53:23 [Debug] JFrog CLI version: 2.52.8
┃ 09:53:23 [Debug] OS/Arch: linux/arm64
┃ 09:53:23 [Debug] Usage info is disabled.
┃ 09:53:23 [Debug] Sending HTTP GET request to: https://artifactory.example.com/artifactory/api/system/ping
┃ 09:53:23 [Debug] Artifactory response: 200 OK
┃ OK
```

Publish a Go library:

```bash
$ export ARTIFACTORY_PASSWORD=xyz
$ dagger call --instance-url=https://artifactory.example.com/artifactory --username=YOUR_USER --password=${ARTIFACTORY_PASSWORD} publishGoLib --repo YOUR_ARTIFACTORY_REPO --src ./testdata --version v0.0.1 --log-level debug
```
