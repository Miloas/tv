# tv
[![Build Status](https://dev.azure.com/0x9357/0x9357/_apis/build/status/Miloas.tv?branchName=master)](https://dev.azure.com/0x9357/0x9357/_build/latest?definitionId=1&branchName=master) [![GitHub release](https://img.shields.io/github/release/Miloas/tv.svg)](https://github.com/Miloas/tv/releases/)

## Installation

Install tv in **macOS** through `brew`

```
brew tap Miloas/tv
brew install tv
```

## Initialize your project

create a `semver.json` file under your project folder.

```json
// [semver.json] you need to setup at least one app
{
  "tv": "0.1.0",
  "app2": "0.1.1"
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
   1.0.15

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
➜  tv git:(master) ./tv patch --help
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

and a git tag `v0.1.1+tv`

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
