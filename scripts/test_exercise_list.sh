#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "exercise_list"

# Setup test data
echo -e "\n${YELLOW}Setting up test data${NC}"
echo "y" | TEST_MODE=true ./bin/tracker exercise add --activity jogging --duration 45 --date 2024-01-08 --notes "Morning run" --completed
echo "y" | TEST_MODE=true ./bin/tracker exercise add --activity cycling --duration 60 --date 2024-01-09 --notes "Evening ride"
echo "y" | TEST_MODE=true ./bin/tracker exercise add --activity other --other-activity "swimming" --duration 30 --date 2024-01-10 --notes "Pool workout" --completed
verify_exercise_file

# Test 1: Basic list with date range
echo -e "\n${YELLOW}Test 1: Basic list${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise list --from 2024-01-01 --to 2024-01-31 2>&1)
assert_output_contains "$output" "Morning run" "Shows first record"
assert_output_contains "$output" "Evening ride" "Shows second record"
assert_output_contains "$output" "Pool workout" "Shows third record"
assert_output_contains "$output" "Total Records    : 3" "Shows correct count"
assert_output_contains "$output" "Total Duration   : 135 minutes" "Shows duration stats"
assert_output_contains "$output" "Completed Records: 2" "Shows completion stats"

# Test 2: Empty range
echo -e "\n${YELLOW}Test 2: Empty range${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise list --from 2023-01-01 --to 2023-12-31 2>&1)
assert_output_contains "$output" "No exercise records found" "Shows empty message"

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