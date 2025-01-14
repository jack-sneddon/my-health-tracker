#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "exercise_get"

# Setup test data
echo -e "\n${YELLOW}Setting up test data${NC}"
# Add a basic exercise record
echo "y" | TEST_MODE=true ./bin/tracker exercise add --activity jogging --duration 45 --date 2024-01-08 --notes "First exercise"
# Add a completed exercise record
echo "y" | TEST_MODE=true ./bin/tracker exercise add --activity cycling --duration 60 --date 2024-01-09 --completed --notes "Evening ride"
# Add an 'other' activity record
echo "y" | TEST_MODE=true ./bin/tracker exercise add --activity other --other-activity "swimming" --duration 30 --date 2024-01-10 --notes "Pool workout"
verify_exercise_file

# Test 1: Get basic record
echo -e "\n${YELLOW}Test 1: Get basic record${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$output" "jogging" "Shows correct activity"
assert_output_contains "$output" "45 minutes" "Shows correct duration"
assert_output_contains "$output" "First exercise" "Shows notes"
assert_output_contains "$output" "Completed:  false" "Shows not completed"

# Test 2: Get completed record
echo -e "\n${YELLOW}Test 2: Get completed record${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-09 2>&1)
assert_output_contains "$output" "cycling" "Shows correct activity"
assert_output_contains "$output" "60 minutes" "Shows correct duration"
assert_output_contains "$output" "Evening ride" "Shows notes"
assert_output_contains "$output" "Completed:  true" "Shows completed status"

# Test 3: Get other activity record
echo -e "\n${YELLOW}Test 3: Get other activity record${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-10 2>&1)
assert_output_contains "$output" "other (swimming)" "Shows other activity type with details"
assert_output_contains "$output" "30 minutes" "Shows correct duration"
assert_output_contains "$output" "Pool workout" "Shows notes"

# Test 4: Get non-existent record
echo -e "\n${YELLOW}Test 4: Get non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-11 2>&1)
assert_output_contains "$output" "not found" "Shows not found message"

# Test 5: Invalid date format
echo -e "\n${YELLOW}Test 5: Invalid date format${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise get --date "invalid" 2>&1)
assert_output_contains "$output" "invalid date format" "Shows invalid date message"

show_test_summary