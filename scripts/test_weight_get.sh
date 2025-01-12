# scripts/test_weight_get.sh
#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "weight_get"

# Setup test data
echo -e "\n${YELLOW}Setting up test data${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
verify_data_file

# Test 1: Get existing record
echo -e "\n${YELLOW}Test 1: Get existing record${NC}"
output=$(TEST_MODE=true ./bin/tracker weight get --date 2024-01-08 2>&1)
assert_output_contains "$output" "185.5" "Shows correct weight"
assert_output_contains "$output" "First weight" "Shows notes"

# Test 2: Get non-existent record
echo -e "\n${YELLOW}Test 2: Get non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker weight get --date 2024-01-09 2>&1)
assert_output_contains "$output" "not found" "Shows not found message"

# Test 3: Invalid date format
echo -e "\n${YELLOW}Test 3: Invalid date format${NC}"
output=$(TEST_MODE=true ./bin/tracker weight get --date "invalid" 2>&1)
assert_output_contains "$output" "invalid date format" "Shows invalid date message"

show_test_summary