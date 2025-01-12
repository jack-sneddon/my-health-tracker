#!/bin/bash

# Create output directory
mkdir -p out

# Color formatting
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "Running all weight command tests..."
echo "=================================="

# Track overall status
FAILED=0

run_test() {
   local test_name=$1
   local test_script=$2

   echo -e "\nRunning $test_name tests..."
   if ./$test_script; then
       echo -e "${GREEN}✓ $test_name tests passed${NC}"
   else
       echo -e "${RED}✗ $test_name tests failed${NC}"
       FAILED=1
   fi
   echo "-----------------------------------"
}

# Run each test script
run_test "weight add" "scripts/test_weight_add.sh"
run_test "weight get" "scripts/test_weight_get.sh"
run_test "weight list" "scripts/test_weight_list.sh"
run_test "weight update" "scripts/test_weight_update.sh"
run_test "weight delete" "scripts/test_weight_delete.sh"

echo -e "\nTest suite completed."
if [ $FAILED -eq 0 ]; then
   echo -e "${GREEN}All test suites passed!${NC}"
   exit 0
else
   echo -e "${RED}Some test suites failed.${NC}"
   exit 1
fi