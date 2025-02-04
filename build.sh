#!/usr/bin/env bash

set -exuo pipefail

mkdir -p bin
go build -o ./bin/server ./cmd/server
