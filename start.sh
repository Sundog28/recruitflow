#!/bin/sh
set -e

echo "Building API..."
go build -o server ./apps/api/cmd/api

echo "Starting API on port $PORT"
exec ./server
