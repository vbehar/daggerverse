# JFrog CLI Dagger Module

This is a [Dagger](https://dagger.io/) module for the [JFrog CLI](https://github.com/jfrog/jfrog-cli/) tool.

Use it to install the JFrog CLI in a Dagger container.

## Usage

### Dagger Shell

To play with the JFrog CLI in a Dagger shell, run:

```bash
$ dagger call install terminal
/ # jf --version
jf version 2.52.8
```

or even with a custom version and base container:

```bash
$ dagger shell --version 2.52.3 install --base alpine
/ # jf --version
jf version 2.52.3
```

### Dagger Go SDK

To use this module from another module:

```go
func doSomething(ctr *Container) {
    ctr = dag.Jfrogcli(JfrogcliOpts{
		Version: "2.52.8",
	}).Install(JfrogcliInstallOpts{
		Base: ctr,
	})

    // ...
}
```
