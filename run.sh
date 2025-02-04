#!/usr/bin/env bash
set -exuo pipefail

./build.sh
DB_CONNECTION=./cookies.sqlite ./bin/server
