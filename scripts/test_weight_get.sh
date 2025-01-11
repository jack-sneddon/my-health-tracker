# scripts/test_weight_get.sh
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

TEST_DATA_DIR=~/.health-tracker/data/test

# Cleanup any existing test data
rm -rf "$TEST_DATA_DIR"

echo -e "${YELLOW}Starting weight get tests...${NC}"

# Setup test data
echo -e "${YELLOW}Setting up test data...${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "Second weight"

# Test successful get - capture both stdout and stderr
echo -e "\n${YELLOW}Test 1: Getting existing record...${NC}"
output=$(TEST_MODE=true ./bin/tracker weight get --date 2024-01-08 2>&1)
echo "$output"

echo -e "\n${YELLOW}Verifying successful get:${NC}"
if [[ "$output" == *"Weight Record"* && "$output" == *"185.5"* ]]; then
    echo -e "${GREEN}✓ Found correct weight${NC}"
else
    echo -e "${RED}✗ Weight not found or incorrect${NC}"
    echo -e "${RED}Output was: $output${NC}"
fi

if [[ "$output" == *"First weight"* ]]; then
    echo -e "${GREEN}✓ Found correct note${NC}"
else
    echo -e "${RED}✗ Note not found or incorrect${NC}"
    echo -e "${RED}Output was: $output${NC}"
fi

# Test non-existent date
echo -e "\n${YELLOW}Test 2: Getting non-existent record...${NC}"
output=$(TEST_MODE=true ./bin/tracker weight get --date 2024-01-07 2>&1)
echo "$output"
if [[ "$output" == *"not found"* ]]; then
    echo -e "${GREEN}✓ Correctly handled non-existent record${NC}"
else
    echo -e "${RED}✗ Failed to handle non-existent record${NC}"
    echo -e "${RED}Output was: $output${NC}"
fi

# Test invalid date format
echo -e "\n${YELLOW}Test 3: Testing invalid date format...${NC}"
output=$(TEST_MODE=true ./bin/tracker weight get --date "invalid-date" 2>&1)
echo "$output"
if [[ "$output" == *"invalid date format"* ]]; then
    echo -e "${GREEN}✓ Correctly handled invalid date format${NC}"
else
    echo -e "${RED}✗ Failed to handle invalid date format${NC}"
    echo -e "${RED}Output was: $output${NC}"
fi

echo -e "\n${YELLOW}Current test data:${NC}"
if [ -f "$TEST_DATA_DIR/weight.json" ]; then
    cat "$TEST_DATA_DIR/weight.json"
else
    echo -e "${RED}No test data file found!${NC}"
fi