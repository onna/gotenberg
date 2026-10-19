#!/bin/bash

set -x

go test -race -covermode=atomic -coverprofile=/tests/coverage.txt ./...
status=$?
go tool cover -html=coverage.txt -o /tests/coverage.html
exit $status
