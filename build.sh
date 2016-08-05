#!/bin/bash

set -e # fail on error

echo "Building porthole"
go test ./...
go install github.com/epond/porthole
echo "Built and tested"
