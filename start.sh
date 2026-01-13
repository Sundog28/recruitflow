#!/usr/bin/env bash
set -euo pipefail

cd apps/api

# build + run
go mod download
go build -o server ./cmd/api
exec ./server
