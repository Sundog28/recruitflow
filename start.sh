#!/bin/sh
set -eu

cd apps/api

go mod download
go build -o server ./cmd/api
exec ./server
