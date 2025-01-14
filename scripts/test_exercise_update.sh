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
output=$(echo "y" | TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --activity cycling 2>&1)
assert_output_contains "$output" "Exercise record updated successfully" "Update succeeded"
record_check=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$record_check" "cycling" "Verifying updated activity"

# Reset data before next test
reset_test_data

# Test 2: Update duration
echo -e "\n${YELLOW}Test 2: Update duration${NC}"
output=$(echo "y" | TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --duration 60 2>&1)
assert_output_contains "$output" "Exercise record updated successfully" "Update succeeded"
record_check=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$record_check" "60" "Verifying updated duration"

# Reset data before next test
reset_test_data

# Test 3: Update to 'other' activity
echo -e "\n${YELLOW}Test 3: Update to other activity${NC}"
output=$(echo "y" | TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --activity other --other-activity "swimming" 2>&1)
assert_output_contains "$output" "Exercise record updated successfully" "Update succeeded"
record_check=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$record_check" "other (swimming)" "Verifying updated activity type"

# Reset data before next test
reset_test_data

# Test 4: Update completion status
echo -e "\n${YELLOW}Test 4: Update completion status${NC}"
output=$(echo "y" | TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --not-completed 2>&1)
assert_output_contains "$output" "Exercise record updated successfully" "Update succeeded"
record_check=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$record_check" "Completed:  false" "Verifying updated completion status"

# Test 5: Invalid activity
echo -e "\n${YELLOW}Test 5: Invalid activity${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --activity invalid 2>&1)
assert_output_contains "$output" "invalid activity type" "Shows invalid activity message"

# Test 6: Other activity without details
echo -e "\n${YELLOW}Test 6: Invalid other activity${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --activity other 2>&1)
assert_output_contains "$output" "other-activity flag is required" "Shows missing other activity message"

# Test 7: Update non-existent record
echo -e "\n${YELLOW}Test 7: Update non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise update --date 2023-01-01 --activity cycling 2>&1)
assert_output_contains "$output" "not found" "Shows not found message"

# Test 8: Duration warning
echo -e "\n${YELLOW}Test 8: Duration warning${NC}"
output=$(echo "n" | TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --duration 120 2>&1)
assert_output_contains "$output" "Duration change" "Shows duration warning"
assert_output_contains "$output" "Operation cancelled" "Shows cancellation message"
record_check=$(TEST_MODE=true ./bin/tracker exercise get --date 2024-01-08 2>&1)
assert_output_contains "$record_check" "45" "Verifying duration unchanged after cancellation"

# Test 9: Conflicting completion flags
echo -e "\n${YELLOW}Test 9: Conflicting completion flags${NC}"
output=$(TEST_MODE=true ./bin/tracker exercise update --date 2024-01-08 --completed --not-completed 2>&1)
assert_output_contains "$output" "cannot use both" "Shows conflicting flags message"

show_test_summary