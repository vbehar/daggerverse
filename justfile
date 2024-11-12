# https://github.com/casey/just

default:
  @just --list

develop-all:
  #!/usr/bin/env sh
  find . -name "dagger.json" -execdir sh -c 'pwd && dagger develop' \;
