# scripts/test_weight_delete.sh
#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "weight_delete"

# Setup test data
echo -e "\n${YELLOW}Setting up test data${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
verify_data_file

# Test 1: Delete existing record
echo -e "\n${YELLOW}Test 1: Delete existing record${NC}"
output=$(echo "y" | TEST_MODE=true ./bin/tracker weight delete w00001 2>&1)
assert_output_contains "$output" "Weight record deleted successfully" "Delete succeeded"

# Test 2: Delete non-existent record
echo -e "\n${YELLOW}Test 2: Delete non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker weight delete w99999 2>&1)
assert_output_contains "$output" "not found" "Shows not found message"

# Test 3: Cancel deletion
echo -e "\n${YELLOW}Test 3: Cancel deletion${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
output=$(echo "n" | TEST_MODE=true ./bin/tracker weight delete w00002 2>&1)
assert_output_contains "$output" "Operation cancelled" "Shows cancellation message"

show_test_summary