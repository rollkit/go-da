# go-da

go-da defines a generic Data Availability interface for modular blockchains.

Note that the rollup clients _do not_ need to depend on any implementation,
they can just make sure that the DA interface is satisfied and start using the
service. This is a key feature of modular blockchains as they can switch the
implementations without having to change the interface.

<!-- markdownlint-disable MD013 -->
[![build-and-test](https://github.com/rollkit/go-da/actions/workflows/ci_release.yml/badge.svg)](https://github.com/rollkit/celestia-da/actions/workflows/ci_release.yml)
[![golangci-lint](https://github.com/rollkit/go-da/actions/workflows/lint.yml/badge.svg)](https://github.com/rollkit/celestia-da/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rollkit/go-da)](https://goreportcard.com/report/github.com/rollkit/celestia-da)
[![codecov](https://codecov.io/gh/rollkit/go-da/branch/main/graph/badge.svg?token=CWGA4RLDS9)](https://codecov.io/gh/rollkit/celestia-da)
[![GoDoc](https://godoc.org/github.com/rollkit/go-da?status.svg)](https://godoc.org/github.com/rollkit/celestia-da)
<!-- markdownlint-enable MD013 -->

## DA Interface

| Method      | Params                        | Return       |
| ----------- |-------------------------------| -------------|
| `MaxBlobSize` |                               | `uint64`       |
| `Get`         | `ids []ID`                      | `[]Blobs`      |
| `GetIDs`      | `height uint64`                 | `[]ID`         |
| `Commit`      | `blobs []Blob`                  | `[]Commitment` |
| `Validate`    | `ids []Blob, proofs []Proof`    | `[]bool`       |

## Implementations

See [celestia-da](https://github.com/rollkit/celestia-da) for the Celestia
implementation.

## Helpful commands

```sh
# Generate protobuf files. Requires docker.
make proto-gen

# Lint protobuf files. Requires docker.
make proto-lint

# Run tests.
make test

# Run linters (requires golangci-lint, markdownlint, hadolint, and yamllint)
make lint
```

## Contributing

We welcome your contributions! Everyone is welcome to contribute, whether it's
in the form of code, documentation, bug reports, feature
requests, or anything else.

If you're looking for issues to work on, try looking at the
[good first issue list](https://github.com/rollkit/go-da/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).
Issues with this tag are suitable for a new external contributor and is a great
way to find something you can help with!

Please join our
[Community Discord](https://discord.com/invite/YsnTPcSfWQ)
to ask questions, discuss your ideas, and connect with other contributors.

## Code of Conduct

See our Code of Conduct [here](https://docs.celestia.org/community/coc).
