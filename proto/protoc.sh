#!/usr/bin/env bash

set -eo pipefail

buf generate --path="./proto/da" --template="buf.gen.yaml" --config="buf.yaml"