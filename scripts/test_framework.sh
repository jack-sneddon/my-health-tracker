# scripts/test_framework.sh
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

# Test data location
TEST_DATA_DIR=~/.health-tracker/data/test

# Test tracking
PASSED=0
FAILED=0
TOTAL=0

# Setup and cleanup
setup_test_env() {
    local test_name=$1
    echo -e "${YELLOW}Starting ${test_name} tests...${NC}"

    echo "Building application..."
    # Format the code
    go fmt ./...
    # Tidy up modules
    go mod tidy
    # Check for errors
    go vet ./...
    # Build the binary
    go build -o bin/tracker ./cmd/tracker

    cleanup_test_data
}

# ... rest of the framework ...

cleanup_test_data() {
    rm -rf "$TEST_DATA_DIR"
}

# Assertions
assert_output_contains() {
    local output=$1
    local expected=$2
    local message=$3
    ((TOTAL++))
    
    if [[ "$output" == *"$expected"* ]]; then
        echo -e "${GREEN}✓ $message${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ $message${NC}"
        echo "Expected: $expected"
        echo "Got: $output"
        ((FAILED++))
    fi  # Fixed the syntax error here
}

# Data verification
verify_data_file() {
    if [ -f "$TEST_DATA_DIR/weight.json" ]; then
        echo -e "${YELLOW}Current test data:${NC}"
        cat "$TEST_DATA_DIR/weight.json"
    else
        echo -e "${RED}No test data file found!${NC}"
        return 1
    fi
}

verify_exercise_file() {
    if [ -f "$TEST_DATA_DIR/exercise.json" ]; then
        echo -e "${YELLOW}Current test data:${NC}"
        cat "$TEST_DATA_DIR/exercise.json"
    else
        echo -e "${RED}No test data file found!${NC}"
        return 1
    fi
}

# Summary
show_test_summary() {
    echo -e "\n${YELLOW}Test Summary${NC}"
    echo "----------------"
    echo "Total Tests: $TOTAL"
    echo "Passed: $PASSED"
    echo "Failed: $FAILED"
    
    if [ $FAILED -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
    else
        echo -e "${RED}Some tests failed.${NC}"
    fi
}