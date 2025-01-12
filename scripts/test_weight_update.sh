# scripts/test_weight_update.sh
#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "weight_update"

# Setup test data
echo -e "\n${YELLOW}Setting up test data${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
verify_data_file

# Test 1: Update weight
echo -e "\n${YELLOW}Test 1: Update weight${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w00001 --value 186.0 2>&1)
assert_output_contains "$output" "Weight record updated successfully" "Update succeeded"
assert_output_contains "$output" "186.0" "Shows updated weight"

# Test 2: Update non-existent record
echo -e "\n${YELLOW}Test 2: Update non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w99999 --value 186.0 2>&1)
assert_output_contains "$output" "not found" "Shows not found message"

# Test 3: Invalid weight value
echo -e "\n${YELLOW}Test 3: Invalid weight value${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w00001 --value 45.0 2>&1)
assert_output_contains "$output" "weight must be between" "Shows invalid weight message"

show_test_summary