#!/bin/bash

# Build script for jit that handles VCS issues
echo "Building jit..."

# Set environment variable to disable VCS stamping
export CGO_ENABLED=0
export GOFLAGS="-buildvcs=false"

# Build and install
go install -buildvcs=false ./cmd/jit

if [ $? -eq 0 ]; then
    echo "Build successful! jit has been installed."
    echo "You can now run: jit --help"
else
    echo "Build failed!"
    exit 1
fi 