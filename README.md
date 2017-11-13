**WARNING:** This project is in alpha

# pmux [![Build Status](https://travis-ci.org/drewschrauf/pmux.svg?branch=master)](https://travis-ci.org/drewschrauf/pmux)

`pmux` is a utility for working with complex stacks spanning multiple repositories and technologies.

## Why?

With some of the more complex systems I've worked on, I've needed to run up to 8 separate terminals with 8 different commands to spin up the full stack. With `pmux`, this can be achieved in one terminal with one command, `pmux run start`.

## Commands

### status

    > pmux status

Show the git status of all projects in the current workspace.

### run

    > pmux run [command]

Run the command (as defined in the config) against all projects in the current workspace.

### exec

    > pmux exec [command]

Run the given, arbitrary command against all projects in the current workspace.

## Config

`pmux` looks for its default config at `~/.pmux/default.yml`. Alternate configs can be used by setting the environment variable `PMUX_WORKSPACE` or by setting the command line flag `--workspace`.

With the following config, `pmux run start` will run the `start` script against projects foo, bar and baz and `pmux run deploy` will run the `deploy` script against project foo.

```yaml
---
projects:
  foo:
    dir: ~/projects/foo
    commands:
      start:
        script: ./start.sh
      deploy:
        script: make deploy
  bar:
    dir: ~/projects/bar
    commands:
      start:
        script: go run project.go
  baz:
    dir: ~/projects/baz
    commands:
      start:
        script: yarn start
        requires:
          - foo
          - bar
```

## Install

Simply get it with `go get`:

```
go get -u github.com/drewschrauf/pmux
```
