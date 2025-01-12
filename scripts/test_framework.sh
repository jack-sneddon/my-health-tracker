# scripts/test_framework.sh
#!/bin/bash

# Colors and formatting
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
NC=$(printf '\033[0m')

# Test data location
TEST_DATA_DIR=~/.health-tracker/data/test
TEST_RESULTS=()
TOTAL_TESTS=0
PASSED_TESTS=0

# Setup functions
setup_test_env() {
    local test_name=$1
    echo -e "${YELLOW}Starting $test_name tests...${NC}"
    cleanup_test_data
}

cleanup_test_data() {
    rm -rf "$TEST_DATA_DIR"
}

create_test_data() {
    echo -e "${YELLOW}Creating test data...${NC}"
    echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
    echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "Second weight"
}

# Test verification functions
assert_success() {
    local description=$1
    local condition=$2
    ((TOTAL_TESTS++))

    if eval "$condition"; then
        echo -e "${GREEN}✓ $description${NC}"
        ((PASSED_TESTS++))
        TEST_RESULTS+=("✓ $description")
    else
        echo -e "${RED}✗ $description${NC}"
        TEST_RESULTS+=("✗ $description")
    fi
}

assert_output_contains() {
    local output=$1
    local expected=$2
    local description=$3
    
    if [[ "$output" == *"$expected"* ]]; then
        echo -e "${GREEN}✓ $description${NC}"
        ((PASSED_TESTS++))
        TEST_RESULTS+=("✓ $description")
    else
        echo -e "${RED}✗ $description${NC}"
        echo -e "${RED}Expected to find: '$expected'${NC}"
        echo -e "${RED}Actual output: '$output'${NC}"
        TEST_RESULTS+=("✗ $description")
    fi
    ((TOTAL_TESTS++))
}

# Test summary
show_test_summary() {
    echo -e "\n${YELLOW}Test Summary${NC}"
    echo "----------------"
    echo "Total Tests: $TOTAL_TESTS"
    echo "Passed: $PASSED_TESTS"
    echo "Failed: $((TOTAL_TESTS - PASSED_TESTS))"

    if [ $TOTAL_TESTS -eq $PASSED_TESTS ]; then
        echo -e "${GREEN}All tests passed!${NC}"
    else
        echo -e "${RED}Some tests failed.${NC}"
    fi
}

# Data verification
verify_data_file() {
    if [ -f "$TEST_DATA_DIR/weight.json" ]; then
        echo -e "\n${YELLOW}Current test data:${NC}"
        cat "$TEST_DATA_DIR/weight.json"
    else
        echo -e "${RED}No test data file found!${NC}"
        return 1
    fi
}

assert_not_contains() {
    local output=$1
    local unexpected=$2
    local description=$3
    
    if [[ "$output" != *"$unexpected"* ]]; then
        echo -e "${GREEN}✓ $description${NC}"
        ((PASSED_TESTS++))
        TEST_RESULTS+=("✓ $description")
    else
        echo -e "${RED}✗ $description${NC}"
        echo -e "${RED}Found unexpected: '$unexpected'${NC}"
        echo -e "${RED}In output: '$output'${NC}"
        TEST_RESULTS+=("✗ $description")
    fi
    ((TOTAL_TESTS++))
}