#!/bin/bash

set -euo pipefail

NEW_MODULE="github.com/sherwin-77/golang-echo"
OLD_MODULE="github.com/sherwin-77/go-echo-template"

# Update the module name in go.mod
go mod edit -module $NEW_MODULE_NAME

# Rename all imported modules in .go files recursively
find . -type f -name "*.go" -not -path "./vendor/*" -exec sed -i -e 's/{OLD_MODULE},{NEW_MODULE}/g' {} \;
