# scripts/test_weight_add.sh
#!/bin/bash

source ./scripts/test_framework.sh

TEST_NAME="weight_add"

# Test cases
test_basic_add() {
    echo -e "\n${YELLOW}Test 1: Basic weight addition${NC}"
    output=$(echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight" 2>&1)

    assert_output_contains "$output" "Weight record added successfully" "Weight record was added"
    assert_output_contains "$output" "185.5" "Correct weight value shown"
    assert_output_contains "$output" "First weight" "Notes were saved"
}

test_duplicate_date() {
    echo -e "\n${YELLOW}Test 2: Duplicate date handling${NC}"
    echo "Current data before duplicate test:"
    verify_data_file

    # Test with 'n' to cancel
    output=$(echo "n" | TEST_MODE=true ./bin/tracker weight add -v 186.0 --date 2024-01-08 2>&1)
    echo "Test output:"
    echo "$output"

    assert_output_contains "$output" "Record already exists for 2024-01-08" "Duplicate date was detected"
    assert_output_contains "$output" "Operation cancelled" "Operation was cancelled"

    # Verify original data wasn't modified
    verify_data=$(cat "$TEST_DATA_DIR/weight.json")
    assert_output_contains "$verify_data" "185.5" "Original weight value remained unchanged"
}

test_invalid_weight() {
    echo -e "\n${YELLOW}Test 3: Invalid weight validation${NC}"
    output=$(TEST_MODE=true ./bin/tracker weight add -v 45.0 --date 2024-01-09 2>&1)
    assert_output_contains "$output" "weight must be between" "Invalid weight was rejected"
}

# Main test execution
main() {
    setup_test_env "$TEST_NAME"

    test_basic_add
    verify_data_file

    test_duplicate_date
    test_invalid_weight

    show_test_summary
}

# Run tests
{
    main
} > "out/${TEST_NAME}_output.txt"

# Show results location
echo "Test results written to out/${TEST_NAME}_output.txt"