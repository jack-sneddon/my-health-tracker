#!/bin/bash

source ./scripts/test_framework.sh

# Initialize test
setup_test_env "exercise_add"

# Test 1: Basic add
echo -e "\n${YELLOW}Test 1: Basic add${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise add --activity jogging --duration 45 --date 2024-01-08 --notes "Morning run" 2>&1)
assert_output_contains "$output" "Exercise record added successfully" "Record was added"
verify_exercise_file

# Test 2: Invalid other activity
echo -e "\n${YELLOW}Test 2: Invalid other activity${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise add --activity other --duration 30 --date 2024-01-09 2>&1)
assert_output_contains "$output" "other-activity flag is required" "Validation failed correctly"

# Test 3: Invalid activity type
echo -e "\n${YELLOW}Test 3: Invalid activity type${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise add --activity invalid --duration 45 --date 2024-01-10 2>&1)
assert_output_contains "$output" "invalid activity type" "Invalid activity rejected"

# Test 4: Invalid duration
echo -e "\n${YELLOW}Test 4: Invalid duration${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise add --activity walking --duration 0 --date 2024-01-11 2>&1)
assert_output_contains "$output" "duration must be greater than 0" "Invalid duration rejected"

# Test 5: Add with completed flag
echo -e "\n${YELLOW}Test 5: Add with completed flag${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise add --activity cycling --duration 60 --date 2024-01-12 --completed 2>&1)
assert_output_contains "$output" "Exercise record added successfully" "Record with completed flag was added"
verify_exercise_file
# Verify specific record content
cycling_record=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-12 2>&1)
assert_output_contains "$cycling_record" "cycling" "Activity type is correct"
assert_output_contains "$cycling_record" "60" "Duration is correct"
assert_output_contains "$cycling_record" "true" "Completed flag is set"

# Test 6: Add with other activity type
echo -e "\n${YELLOW}Test 6: Add other activity${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise add --activity other --other-activity "swimming" --duration 30 --date 2024-01-13 2>&1)
assert_output_contains "$output" "Exercise record added successfully" "Other activity record was added"
verify_exercise_file
# Verify specific record content
swimming_record=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-13 2>&1)
assert_output_contains "$swimming_record" "other" "Activity type is correct"
assert_output_contains "$swimming_record" "swimming" "Other activity is specified"
assert_output_contains "$swimming_record" "30" "Duration is correct"

# Show results
show_test_summary