[![Unit tests](https://github.com/archspec/archspec-go/actions/workflows/unit_tests.yml/badge.svg)](https://github.com/archspec/archspec-go/actions/workflows/unit_tests.yml)
[![codecov](https://codecov.io/gh/archspec/archspec-go/branch/master/graph/badge.svg)](https://codecov.io/gh/archspec/archspec-go)

# Archspec (Go bindings)

Archspec aims at providing a standard set of human-understandable labels for
various aspects of a system architecture  like CPU, network fabrics, etc. and
APIs to detect, query and compare them.

This project grew out of [Spack](https://spack.io/) and is currently under
active development. At present it supports APIs to detect and model
compatibility relationships among different CPU microarchitectures.

## Development

You will need Go minimum version 1.19. From the root of the repository, you'll want to ensure you 
have the git submodule init and updated:

```bash
git submodule init
git submodule update
```

This adds the json microarchitecture files in [archspec/json](archspec/json). To test:

```bash
go test ./...
```

## License

Archspec is distributed under the terms of both the MIT license and the
Apache License (Version 2.0). Users may choose either license, at their
option.

All new contributions must be made under both the MIT and Apache-2.0
licenses.

See [LICENSE-MIT](https://github.com/archspec/archspec-go/blob/master/LICENSE-MIT),
[LICENSE-APACHE](https://github.com/archspec/archspec-go/blob/master/LICENSE-APACHE),
[COPYRIGHT](https://github.com/archspec/archspec-go/blob/master/COPYRIGHT), and
[NOTICE](https://github.com/archspec/archspec-go/blob/master/NOTICE) for details.

SPDX-License-Identifier: (Apache-2.0 OR MIT)

LLNL-CODE-811653
