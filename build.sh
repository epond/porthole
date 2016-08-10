#!/bin/bash

set -e # fail on error

echo "Building porthole"
go get github.com/djherbis/times
go test ./...
go install github.com/epond/porthole
echo "Built and tested"
