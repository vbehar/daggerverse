# JFrog CLI Dagger Module

This is a [Dagger](https://dagger.io/) module for the [JFrog CLI](https://github.com/jfrog/jfrog-cli/) tool.

Use it to install the JFrog CLI in a Dagger container.

## Installation

```bash
$ dagger install github.com/vbehar/daggerverse/jfrogcli
```

## Usage

### Shell

Test the output of the `jf --version` command:

```bash
$ dagger call -i -m github.com/vbehar/daggerverse/jfrogcli install with-exec --args jf,--version stdout
```

### Dagger Go SDK

To use this module from another module:

```go
func doSomething(ctr *Container) {
    ctr = dag.Jfrogcli(dagger.JfrogcliOpts{
		// Version: "2.71.0", // optional
	}).Install(dagger.JfrogcliInstallOpts{
		Base: ctr,
	})

    // ...
}
```
