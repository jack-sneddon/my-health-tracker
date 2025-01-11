# scripts/test_weight_list.sh
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Setup test directory
TEST_DATA_DIR=~/.health-tracker/data/test

# Cleanup and setup
rm -rf "$TEST_DATA_DIR"

echo -e "${YELLOW}Starting weight list tests...${NC}"

# Setup test data
echo -e "${YELLOW}Setting up test data...${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "Second weight"

echo -e "${YELLOW}Testing list with date range...${NC}"
output=$(TEST_MODE=true ./bin/tracker weight list --from 2024-01-01 --to 2024-01-31)
echo "$output"

echo -e "\n${YELLOW}Verifying output:${NC}"
if [[ "$output" == *"185.5"* && "$output" == *"184.8"* ]]; then
    echo -e "${GREEN}✓ Found both weight records${NC}"
else
    echo -e "${RED}✗ Missing weight records${NC}"
fi

if [[ "$output" == *"First weight"* && "$output" == *"Second weight"* ]]; then
    echo -e "${GREEN}✓ Found both notes${NC}"
else
    echo -e "${RED}✗ Missing notes${NC}"
fi

# Updated to match actual output format
if [[ "$output" == *"Total Records : 2"* ]]; then
    echo -e "${GREEN}✓ Record count correct${NC}"
else
    echo -e "${RED}✗ Record count incorrect or missing${NC}"
fi

echo -e "\n${YELLOW}Current test data:${NC}"
if [ -f "$TEST_DATA_DIR/weight.json" ]; then
    cat "$TEST_DATA_DIR/weight.json"
else
    echo -e "${RED}No test data file found!${NC}"
fi