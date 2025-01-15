#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "exercise_update"

# Helper function to reset test data
reset_test_data() {
    cleanup_test_data
    echo "y" | TEST_MODE=true ./bin/tracker exercise add --activity jogging --duration 45 --date 2024-01-08 --notes "First exercise" --completed
    verify_exercise_file
}

# Initial setup
echo -e "\n${YELLOW}Setting up test data${NC}"
reset_test_data

# Test 1: Update activity
echo -e "\n${YELLOW}Test 1: Update activity${NC}"
echo -e "\n${YELLOW}Before update:${NC}"
TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08
output=$(echo "y" | TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --activity cycling 2>&1)
assert_output_contains "$output" "Exercise record updated successfully" "Update succeeded"
echo -e "\n${YELLOW}After update:${NC}"
verify_exercise_file
record_check=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$record_check" "cycling" "Verifying updated activity"

# Reset data before next test
reset_test_data

# Test 2: Update duration
echo -e "\n${YELLOW}Test 2: Update duration${NC}"
echo -e "\n${YELLOW}Before update:${NC}"
TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08
output=$(echo "y" | TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --duration 60 2>&1)
assert_output_contains "$output" "Exercise record updated successfully" "Update succeeded"
echo -e "\n${YELLOW}After update:${NC}"
verify_exercise_file
record_check=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$record_check" "60" "Verifying updated duration"

# Test 3: Single day
echo -e "\n${YELLOW}Test 3: Single day${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise list --from 2024-01-08 --to 2024-01-08 2>&1)
assert_output_contains "$output" "Total Records    : 1" "Shows single record"
assert_output_contains "$output" "Morning run" "Shows correct record"

# Test 4: Week flag
echo -e "\n${YELLOW}Test 4: Week flag${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise list --from 2024-01-08 --to 2024-01-14 2>&1)
assert_output_contains "$output" "Morning run" "Shows records within week"
assert_output_contains "$output" "Evening ride" "Shows all week records"
assert_output_contains "$output" "Pool workout" "Shows all week records"

# Test 5: Month flag
echo -e "\n${YELLOW}Test 5: Month flag${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise list --month 2>&1)
assert_output_contains "$output" "Exercise Records from 2024-01-01 to 2024-01-31" "Shows correct month range"
assert_output_contains "$output" "Total Records    : 3" "Shows all records for month"

# Test 6: Default range (no dates specified)
echo -e "\n${YELLOW}Test 6: Default range${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise list 2>&1)
assert_output_contains "$output" "Total Records" "Shows records for default range"

# Test 7: Invalid date range
echo -e "\n${YELLOW}Test 7: Invalid date range${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise list --from 2024-01-31 --to 2024-01-01 2>&1)
assert_output_contains "$output" "'from' date must be before 'to' date" "Shows invalid range message"

# Show results
show_test_summary