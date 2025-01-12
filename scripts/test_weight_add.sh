# scripts/test_weight_add.sh
#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "weight_add"

# Test 1: Basic add
echo -e "\n${YELLOW}Test 1: Basic add${NC}"
output=$(echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight" 2>&1)
assert_output_contains "$output" "Weight record added successfully" "Record was added"
verify_data_file

# Test 2: Duplicate date
echo -e "\n${YELLOW}Test 2: Duplicate date${NC}"
output=$(echo "n" | TEST_MODE=true ./bin/tracker weight add -v 186.0 --date 2024-01-08 2>&1)
assert_output_contains "$output" "Record already exists" "Duplicate detected"
assert_output_contains "$output" "Operation cancelled" "Operation cancelled"

# Test 3: Invalid weight
echo -e "\n${YELLOW}Test 3: Invalid weight${NC}"
output=$(TEST_MODE=true ./bin/tracker weight add -v 45.0 --date 2024-01-09 2>&1)
assert_output_contains "$output" "weight must be between" "Invalid weight rejected"

# Show results
show_test_summary