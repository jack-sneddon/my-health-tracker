#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "exercise_delete"

# Helper function to reset test data
reset_test_data() {
    cleanup_test_data
    echo "y" | TEST_MODE=true ./bin/tracker exercise add --activity jogging --duration 45 --date 2024-01-08 --notes "First exercise" --completed
    verify_exercise_file
}

# Initial setup
echo -e "\n${YELLOW}Setting up test data${NC}"
reset_test_data

# Test 1: Delete existing record
echo -e "\n${YELLOW}Test 1: Delete existing record${NC}"
output=$(echo "y" | TEST_MODE=true ./bin/tracker exercise delete --date 2024-01-08 2>&1)
assert_output_contains "$output" "Exercise record deleted successfully" "Delete succeeded"
verify_exercise_file

# Reset data before next test
reset_test_data

# Test 2: Delete cancelled
echo -e "\n${YELLOW}Test 2: Delete cancelled${NC}"
output=$(echo "n" | TEST_MODE=true ./bin/tracker exercise delete --date 2024-01-08 2>&1)
assert_output_contains "$output" "Operation cancelled" "Shows cancellation message"
# Verify record still exists
record_check=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$record_check" "jogging" "Record still exists"

# Test 3: Delete non-existent record
echo -e "\n${YELLOW}Test 3: Delete non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise delete --date 2023-01-01 2>&1)
assert_output_contains "$output" "not found" "Shows not found message"

# Test 4: Invalid date format
echo -e "\n${YELLOW}Test 4: Invalid date format${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise delete --date "invalid" 2>&1)
assert_output_contains "$output" "invalid date format" "Shows invalid date message"

show_test_summary