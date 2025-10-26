#!/bin/bash

# This script reorganizes Go test files into a single 'tests' directory.

TEST_DIR="tests"

echo "--- Creating test directory: $TEST_DIR ---"
mkdir -p "$TEST_DIR"

echo "--- Moving and renaming test files ---"

# Move the storage test
if [ -f "internal/storage/sqlite_test.go" ]; then
    echo "Moving internal/storage/sqlite_test.go -> tests/storage_test.go"
    mv internal/storage/sqlite_test.go tests/storage_test.go
fi

# The old main_test.go is obsolete with the new architecture, so we'll remove it.
if [ -f "main_test.go" ]; then
    echo "Removing obsolete main_test.go"
    rm main_test.go
fi

echo "--- Test file organization complete. ---"
