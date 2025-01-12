# scripts/test_weight_list.sh
#!/bin/bash

source ./scripts/test_framework.sh

TEST_NAME="weight_list"

# Add helper functions for date handling
get_test_dates() {
    # Instead of using current year, we'll set a fixed test date context
    # This ensures tests work regardless of when they're run
    BASE_TEST_YEAR="2024"
    BASE_TEST_MONTH="01"
    BASE_TEST_DAY="15"

    # Return common test dates
    echo "${BASE_TEST_YEAR}-${BASE_TEST_MONTH}-${BASE_TEST_DAY}"
}

setup_test_data() {
    echo -e "\n${YELLOW}Setting up test data...${NC}"
    base_date=$(get_test_dates)

    # Create two records: day before and base date
    day_before=$(date -j -v-1d -f "%Y-%m-%d" "$base_date" +%Y-%m-%d)

    echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date "$day_before" --notes "First weight"
    echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date "$base_date" --notes "Second weight"

    verify_data_file
}

test_basic_list() {
    echo -e "\n${YELLOW}Test 1: Basic list with date range${NC}"
    base_date=$(get_test_dates)
    year=${base_date%%-*}

    output=$(TEST_MODE=true ./bin/tracker weight list --from "${year}-01-01" --to "${year}-01-31" 2>&1)
    echo "Test output:"
    echo "$output"

    assert_output_contains "$output" "185.5" "First weight record shown"
    assert_output_contains "$output" "184.8" "Second weight record shown"
    assert_output_contains "$output" "First weight" "Notes are displayed"
    assert_output_contains "$output" "Total Records : 2" "Correct record count shown"
    assert_output_contains "$output" "Average Weight: 185.2 lbs" "Average weight calculated correctly"
    assert_output_contains "$output" "Weight Range  : 184.8 - 185.5 lbs" "Weight range shown correctly"
}

test_empty_date_range() {
    echo -e "\n${YELLOW}Test 2: Empty date range${NC}"
    output=$(TEST_MODE=true ./bin/tracker weight list --from 2023-01-01 --to 2023-12-31 2>&1)
    echo "Test output:"
    echo "$output"

    assert_output_contains "$output" "No weight records found between 2023-01-01 and 2023-12-31" "Empty result set handled correctly"
}

test_invalid_date_range() {
    echo -e "\n${YELLOW}Test 3: Invalid date range${NC}"
    output=$(TEST_MODE=true ./bin/tracker weight list --from 2024-01-31 --to 2024-01-01 2>&1)
    echo "Test output:"
    echo "$output"

    assert_output_contains "$output" "'from' date must be before 'to' date" "Invalid date range detected"
}

test_default_range() {
    echo -e "\n${YELLOW}Test 4: Default date range (no dates specified)${NC}"
    output=$(TEST_MODE=true ./bin/tracker weight list 2>&1)
    echo "Test output:"
    echo "$output"

    assert_output_contains "$output" "185.5" "Records shown in default range"
    assert_output_contains "$output" "Total Records" "Summary shown for default range"
}

test_single_day_range() {
    echo -e "\n${YELLOW}Test 5: Single day range${NC}"
    today=$(date +%Y-%m-%d)

    output=$(TEST_MODE=true ./bin/tracker weight list --from "$today" --to "$today" 2>&1)
    echo "Test output:"
    echo "$output"

    assert_output_contains "$output" "184.8" "Single day record shown"
    assert_output_contains "$output" "Total Records : 1" "Single record count correct"
}

test_invalid_date_format() {
    echo -e "\n${YELLOW}Test 6: Invalid date format${NC}"
    output=$(TEST_MODE=true ./bin/tracker weight list --from "invalid" --to 2024-01-31 2>&1)
    echo "Test output:"
    echo "$output"

    assert_output_contains "$output" "invalid date format" "Invalid date format detected"
}

test_last_week() {
    echo -e "\n${YELLOW}Test 7: Last week range${NC}"
    output=$(TEST_MODE=true ./bin/tracker weight list --week 2>&1)
    echo "Test output:"
    echo "$output"

    # These records should be within the last week
    assert_output_contains "$output" "Weight Records from" "Date range header shown"
    assert_output_contains "$output" "Total Records : 2" "Shows correct number of records"
    assert_output_contains "$output" "185.5" "Shows first weight"
    assert_output_contains "$output" "184.8" "Shows second weight"
}

test_last_month() {
    echo -e "\n${YELLOW}Test 8: Last month range${NC}"
    output=$(TEST_MODE=true ./bin/tracker weight list --month 2>&1)
    echo "Test output:"
    echo "$output"

    # These records should be within the last month
    assert_output_contains "$output" "Weight Records from" "Date range header shown"
    assert_output_contains "$output" "Total Records : 2" "Shows correct number of records"
    assert_output_contains "$output" "185.5" "Shows first weight"
    assert_output_contains "$output" "184.8" "Shows second weight"
}

test_quick_range_precedence() {
    echo -e "\n${YELLOW}Test 9: Quick range takes precedence over from/to${NC}"
    # Add a test record from last year
    last_year=$(date -v-1y +%Y-%m-%d)
    echo "y" | TEST_MODE=true ./bin/tracker weight add -v 183.0 --date "$last_year" --notes "Old weight"

    # When using --week, should only show recent records even if from/to specifies a wider range
    output=$(TEST_MODE=true ./bin/tracker weight list --week --from 2023-01-01 --to 2025-12-31 2>&1)
    echo "Test output:"
    echo "$output"

    assert_output_contains "$output" "Total Records : 2" "Shows only records within last week"
    assert_not_contains "$output" "183.0" "Does not show old record"
}

# Setup test data
setup_test_data() {
    echo -e "\n${YELLOW}Setting up test data...${NC}"
    # Use current year for test data
    current_year=$(date +%Y)
    today=$(date +%Y-%m-%d)
    yesterday=$(date -v-1d +%Y-%m-%d)

    echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date "$yesterday" --notes "First weight"
    echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date "$today" --notes "Second weight"

    verify_data_file
}

# Main test execution
main() {
    setup_test_env "$TEST_NAME"
    setup_test_data

    test_basic_list
    test_empty_date_range
    test_invalid_date_range
    test_default_range
    test_single_day_range
    test_invalid_date_format
    test_last_week
    test_last_month
    test_quick_range_precedence

    show_test_summary
}

# Run tests
{
    main
} > "out/${TEST_NAME}_output.txt"

# Show results location
echo "Test results written to out/${TEST_NAME}_output.txt"