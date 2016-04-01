# namerctl #

A utility for controlling [namerd](https://github.com/BuoyantIO/linkerd/tree/master/namerd).

This utility _will change_ drastically in the near future.

## Installation ##

```
:; go get -u github.com/BuoyantIO/namerctl
:; go install github.com/BuoyantIO/namerctl
```

## Usage ##

```
:; namerctl dtab help
Control namer's dtab interface

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

Use "namerctl dtab [command] --help" for more information about a
command.
```

## License ##

Copyright 2016, Buoyant Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
these files except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
