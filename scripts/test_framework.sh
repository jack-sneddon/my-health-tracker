# scripts/test_framework.sh
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test data location
TEST_DATA_DIR=~/.health-tracker/data/test

# Cleanup function
cleanup() {
    if [ -d "$TEST_DATA_DIR" ]; then
        echo -e "${YELLOW}Cleaning up test data directory: $TEST_DATA_DIR${NC}"
        rm -rf "$TEST_DATA_DIR"
    fi
}

# Setup function
setup() {
    cleanup
    echo -e "${YELLOW}Formatting code...${NC}"
    go fmt ./...
    
    echo -e "${YELLOW}Building application...${NC}"
    go build -o bin/tracker ./cmd/tracker
    
    # Verify test directory creation
    if [ ! -d "$TEST_DATA_DIR" ]; then
        echo -e "${YELLOW}Test directory will be created on first command${NC}"
    fi
}

# Run test command with test mode flag
run_test_cmd() {
    local cmd="$1"
    TEST_MODE=true $cmd
}

# Run test and compare output
assert_output() {
    local cmd="$1"
    local expected="$2"
    local description="$3"

    echo "Testing: $description"
    local output=$(TEST_MODE=true $cmd)
    
    if [[ "$output" == *"$expected"* ]]; then
        echo -e "${GREEN}✓ Test passed: $description${NC}"
        
        # Verify test directory was created
        if [ -d "$TEST_DATA_DIR" ]; then
            echo -e "${GREEN}✓ Test directory exists: $TEST_DATA_DIR${NC}"
        else
            echo -e "${RED}✗ Test directory was not created${NC}"
            return 1
        fi
        
        return 0
    else
        echo -e "${RED}✗ Test failed: $description${NC}"
        echo "Expected to contain: $expected"
        echo "Got: $output"
        return 1
    fi
}

# Verify directory structure
verify_test_structure() {
    if [ ! -d "$TEST_DATA_DIR" ]; then
        echo -e "${RED}✗ Test directory not found: $TEST_DATA_DIR${NC}"
        return 1
    fi
    
    for file in weight.json exercise.json fasting.json soda.json; do
        if [ ! -f "$TEST_DATA_DIR/$file" ]; then
            echo -e "${RED}✗ Test file not found: $TEST_DATA_DIR/$file${NC}"
            return 1
        fi
    done
    
    echo -e "${GREEN}✓ Test directory structure verified${NC}"
    return 0
}