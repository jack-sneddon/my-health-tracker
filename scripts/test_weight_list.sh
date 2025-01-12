# scripts/test_weight_list.sh
#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "weight_list"

# Setup test data
echo -e "\n${YELLOW}Setting up test data${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "Second weight"
verify_data_file

# Test 1: Basic list
echo -e "\n${YELLOW}Test 1: Basic list${NC}"
output=$(TEST_MODE=true ./bin/tracker weight list --from 2024-01-01 --to 2024-01-31 2>&1)
assert_output_contains "$output" "First weight" "Shows first record"
assert_output_contains "$output" "Second weight" "Shows second record"
assert_output_contains "$output" "Total Records : 2" "Shows correct count"
assert_output_contains "$output" "Average Weight: 185.2" "Shows average"

# Test 2: Empty range
echo -e "\n${YELLOW}Test 2: Empty range${NC}"
output=$(TEST_MODE=true ./bin/tracker weight list --from 2023-01-01 --to 2023-12-31 2>&1)
assert_output_contains "$output" "No weight records found" "Shows empty message"

# Test 3: Single day
echo -e "\n${YELLOW}Test 3: Single day${NC}"
output=$(TEST_MODE=true ./bin/tracker weight list --from 2024-01-08 --to 2024-01-08 2>&1)
assert_output_contains "$output" "Total Records : 1" "Shows single record"
assert_output_contains "$output" "185.5" "Shows correct weight"

# Show results
show_test_summary