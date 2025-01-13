# scripts/test_weight_update.sh
#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "weight_update"

# Helper function to reset test data
reset_test_data() {
   cleanup_test_data
   echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
   verify_data_file
}

# Initial setup
echo -e "\n${YELLOW}Setting up test data${NC}"
reset_test_data

# Test 1: Update weight
echo -e "\n${YELLOW}Test 1: Update weight${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w00001 --value 186.0 2>&1)
assert_output_contains "$output" "Weight record updated successfully" "Update succeeded"
assert_output_contains "$output" "186.0" "Shows updated weight"

# Reset data before next test
reset_test_data

# Test 2: Update non-existent record
echo -e "\n${YELLOW}Test 2: Update non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w99999 --value 186.0 2>&1)
assert_output_contains "$output" "not found" "Shows not found message"

# Reset data before next test
reset_test_data

# Test 3: Invalid weight value
echo -e "\n${YELLOW}Test 3: Invalid weight value${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w00001 --value 45.0 2>&1)
assert_output_contains "$output" "weight must be between" "Shows invalid weight message"

# Reset data before next test
reset_test_data

# Test 4: Update notes only
echo -e "\n${YELLOW}Test 4: Update notes only${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w00001 --notes "Updated notes" 2>&1)
assert_output_contains "$output" "Updated notes" "Shows updated notes"

# Reset data before next test
reset_test_data

# Test 5: Update with warning (significant weight change)
echo -e "\n${YELLOW}Test 5: Update with warning${NC}"
# Add debug output
echo "Current weight before update:"
TEST_MODE=true ./bin/tracker weight get --date 2024-01-08

output=$(echo "n" | TEST_MODE=true ./bin/tracker weight update w00001 --value 200.0 2>&1)
echo "Update command output:"
echo "$output"

assert_output_contains "$output" "seems unusual" "Shows warning message"
assert_output_contains "$output" "Operation cancelled" "Shows cancellation message"

# Verify weight wasn't changed
verify_output=$(TEST_MODE=true ./bin/tracker weight get --date 2024-01-08 2>&1)
assert_output_contains "$verify_output" "185.5" "Weight remained unchanged after cancellation"

# Reset data before next test
reset_test_data

# Test 6: Update both fields
echo -e "\n${YELLOW}Test 6: Update both fields${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w00001 --value 186.0 --notes "Both updated" 2>&1)
assert_output_contains "$output" "186.0" "Shows updated weight"
assert_output_contains "$output" "Both updated" "Shows updated notes"

show_test_summary