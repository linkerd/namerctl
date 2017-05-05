# namerctl #

[![GitHub license][license-badge]](LICENSE)
[![Circle CI][namerctl-ci-badge]][namerctl-ci]
[![Docker Pulls][docker-badge]][docker]

A utility for controlling [namerd](https://github.com/linkerd/linkerd/tree/master/namerd).

This utility _will change_ drastically in the near future.

## Installation ##

```
:; go get -u github.com/linkerd/namerctl
:; go install github.com/linkerd/namerctl
```

## Usage ##

```
$ namerctl help
namerd manages delegation tables for linkerd.

namerctl looks for a configuration file in the current working
directory or any of its parent directories. Configuration files are
named .namerctl.<ext> where <ext> is describes one of several formats
including yaml, json, toml, etc.  "base-url" is currently the only
supported configuration.  Furthermore, the base url may be specified
via the NAMERCTL_BASE_URL environment variable.

Find more information at https://linkerd.io

Usage:
  namerctl [command]

Available Commands:
  dtab        Control namerd's delegation tables

Flags:
      --base-url string   namer location (e.g. http://namerd.example.com:4080)
      --config string     config file

Use "namerctl [command] --help" for more information about a command.
```
```
$ namerctl dtab help
Control namerd's delegation tables

Usage:
  namerctl dtab [command]

Available Commands:
  list        List delegation table names
  get         Get a delegation table by name
  create      Create a new delegation table.
  update      Update a delegation table.
  delete      Delete a delegation by name.

Global Flags:
      --base-url string   namer location (e.g. http://namerd.example.com:4080)
      --config string     config file

Use "namerctl dtab [command] --help" for more information about a command.
```

## Development ##

This project uses [godep](https://github.com/tools/godep) for managing go
dependencies in the `vendor/` directory. To update dependencies, run:

```
$ go get github.com/tools/godep
$ go install github.com/tools/godep
$ godep save
```

## Docker ##

### Running ###

To use the [public image][docker], run:

```
$ docker run --rm linkerd/namerctl:latest --help
namerd manages delegation tables for linkerd.
```

### Building ###

To build your own copy of the image from source, run:

```
$ docker build -t linkerd/namerctl:latest .
```

## License ##

Copyright 2016-2017, Buoyant Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
these files except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.

<!-- references -->
[license-badge]: https://img.shields.io/github/license/linkerd/namerctl.svg
[namerctl-ci-badge]: https://circleci.com/gh/linkerd/namerctl/tree/master.svg?style=shield
[namerctl-ci]: https://circleci.com/gh/linkerd/namerctl
[docker-badge]: https://img.shields.io/docker/pulls/linkerd/namerctl.svg
[docker]: https://hub.docker.com/r/linkerd/namerctl/
