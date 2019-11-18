# tv
[![Build Status](https://github.com/Miloas/tv/workflows/build/badge.svg)](https://github.com/Miloas/tv/workflows/build/badge.svg)
[![GitHub release](https://img.shields.io/github/release/Miloas/tv.svg)](https://github.com/Miloas/tv/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/Miloas/tv)](https://goreportcard.com/report/github.com/Miloas/tv)

## Installation

Install tv in **macOS** through `brew`

```
brew tap Miloas/tv
brew install tv
```

## Initialize your project

Create a `semver.json` file using `tv init app app2` under your project folder. You need to setup at least one app.

```json
{
  "app": "0.0.0",
  "app2": "0.0.0"
}
```

## Usage

```
➜  tv git:(master) tv help
NAME:
   tv - tag version for ur f** awesome project

USAGE:
   tv [global options] command [command options] [arguments...]

VERSION:
   1.1.0

COMMANDS:
   patch       patch version, v0.0.1 -> v0.0.2
   major       major version, v0.0.1 -> v1.0.1
   minor       minor version, v0.0.1 -> v0.1.1
   prerelease  prerelease version, v0.0.1-alpha.1 -> v0.0.1-alpha.2
   version     set specific version
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

Options for commands:

```
➜  tv git:(master) tv patch --help
NAME:
   tv patch - patch version, v0.0.1 -> v0.0.2

USAGE:
   tv patch [command options] [arguments...]

OPTIONS:
   --target value, -t value  set target app
   --pure, -p                create tag without app name
   --all, -a                 upgrade version of all apps
   --dry-run                 do a fake action, won't create real tag
```

## Examples

For such a `semver.json` file

```json
{
  "tv": "0.1.0",
  "app": "0.1.10"
}
```

---

Run `tv patch`, and you will get:

```json
{
  "tv": "0.1.1",
  "app": "0.1.10"
}
```

and a git tag `v0.1.1+tv`. `tv` is selected if you didn't pass a target to tv

---

Run `tv patch -t app`, and you will get:

```json
{
  "tv": "0.1.0",
  "app": "0.1.11"
}
```

and a git tag `v0.1.11+app`

---

Run `tv minor -a`, and you will get:

```json
{
  "tv": "0.2.0",
  "app": "0.2.0"
}
```

and two git tags [`v0.2.0+tv`, `v0.2.0+app`]

---

Run `tv major -ap`, and you will get:

```json
{
  "tv": "1.0.0",
  "app": "1.0.0"
}
```

and only one git tag `v1.0.0`
