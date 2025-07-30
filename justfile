# https://github.com/casey/just

default:
  @just --list

develop-all:
  #!/usr/bin/env sh
  find . -name "dagger.json" -execdir sh -c 'pwd && dagger develop' \;

develop-compat-all:
  #!/usr/bin/env sh
  find . -name "dagger.json" -execdir sh -c 'pwd && dagger develop --compat=$(jq -r '.engineVersion' dagger.json)' \;
