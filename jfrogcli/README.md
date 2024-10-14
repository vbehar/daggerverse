# JFrog CLI Dagger Module

This is a [Dagger](https://dagger.io/) module for the [JFrog CLI](https://github.com/jfrog/jfrog-cli/) tool.

Use it to install the JFrog CLI in a Dagger container.

Read the documentation at <https://daggerverse.dev/mod/github.com/vbehar/daggerverse/jfrogcli>.

## Usage

Test the output of the `jf --version` command:

```bash
$ dagger call -i -m github.com/vbehar/daggerverse/jfrogcli \
	install with-exec --args jf,--version \
	stdout
```
