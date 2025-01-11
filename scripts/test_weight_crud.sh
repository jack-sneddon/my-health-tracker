#!/bin/bash
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