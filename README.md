# cloudping

[![CircleCI](https://img.shields.io/circleci/project/github/badges/shields/master.svg?style=for-the-badge)](https://circleci.com/gh/estahn/cloudping)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/estahn/cloudping)
[![Github All Releases](https://img.shields.io/github/downloads/estahn/cloudping/total.svg?style=for-the-badge)](https://github.com/estahn/cloudping/releases)

cloudping identifies the cloud provider regions geographically closest and returns them in order of lowest to highest latency.

> Inspired by [CloudPing.info](https://www.cloudping.info/).

## Table of Contents

<!-- toc -->

- [Usage](#usage)
- [Installation](#installation)
  * [Binaries](#binaries)
  * [Via Go](#via-go)
  * [Homebrew on macOS](#homebrew-on-macos)
- [Why?](#why)
- [Similar projects](#similar-projects)
- [Contributing](#contributing)

<!-- tocstop -->

## Usage

```console
$ cloudping -h
cloudping identifies the cloud provider regions geographically closest
and returns them in order of lowest to highest latency.

Usage:
  cloudping [flags]
  cloudping [command]

Available Commands:
  help        Help about any command
  version     Print the version number of cloudping

Flags:
  -h, --help              help for cloudping
      --limit int         Limits the number of regions returned
      --output string     Output format. One of: txt, json, yaml (default "txt")
      --provider string   Cloud provider (default "aws")
      --regions strings   Limits checks to specific regions
      --timeout int       Timeout for each region in milliseconds (default 500)

Use "cloudping [command] --help" for more information about a command.
```

The following example checks and returns only 2 regions:

```console
$ cloudping --provider=aws --regions=ap-southeast-2,us-east-1
ap-southeast-2
us-east-1
```

## Installation

Here are a few methods to install `cloudping`.

### Binaries

```console
$ curl -sSL https://raw.githubusercontent.com/estahn/cloudping/master/godownloader.sh | sh
```

### Via Go

```console
$ go get github.com/estahn/cloudping
```

### Homebrew on macOS

If you are on macOS and using [Homebrew](https://brew.sh/) package manager, you can install cloudping with Homebrew.

1. Run the installation command:
   ```console
   $ brew install estahn/cloudping
   ```
2. Test to ensure the version you installed is sufficiently up-to-date:
   ```console
   $ cloudping version
   ```

## Why? 

The idea came from the need to download images from the geographically closest docker registry.
We operate our Kubernetes cluster in Sydney/Australia but use CircleCI operating in the US.
Because AWS ECR doesn't provide a common endpoint with geographically distributed backend we push our images to both locations.
Within our Makefile we can use `cloudping` to identify if images should be pulled from the US or Sydney. 

## Similar projects

* [CloudPing.info](https://www.cloudping.info/)
* [AWS Inter-Region Latency Monitoring](https://www.cloudping.co/)
* [GCP ping](http://www.gcping.com/)

## Contributing

Contributions are greatly appreciated.
The maintainers actively manage the issues list, and try to highlight issues suitable for newcomers.
The project follows the typical GitHub pull request model.
See " [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/) " for more details.
Before starting any work, please either comment on an existing issue, or file a new one.
