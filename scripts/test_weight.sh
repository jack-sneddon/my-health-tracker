#!/bin/bash

# Create output directory
mkdir -p out

# Run all tests and capture both stdout and stderr
{
    echo "Running weight command tests..."
    echo "================================"
    echo

    # First add test data
    ./scripts/test_weight_add.sh

    # Test get command
    ./scripts/test_weight_get.sh

    # Test list command
    ./scripts/test_weight_list.sh

    # Test update command
    ./scripts/test_weight_update.sh

    # Test delete command
    ./scripts/test_weight_delete.sh

} > out/test_weight.txt 2>&1

echo "Tests complete. Results written to out/test_weight.txt"