# scripts/test_weight_add.sh
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Setup
TEST_DATA_DIR=~/.health-tracker/data/test

echo -e "${YELLOW}Starting weight add tests...${NC}"

# Cleanup any existing test data
if [ -d "$TEST_DATA_DIR" ]; then
    echo -e "${YELLOW}Cleaning up previous test data${NC}"
    rm -rf "$TEST_DATA_DIR"
fi

echo -e "${YELLOW}Adding first weight record...${NC}"
output=$(echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight")
echo "$output"

echo -e "\n${YELLOW}Adding second weight record...${NC}"
output=$(echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "Second weight")
echo "$output"

echo -e "\n${YELLOW}Verifying data storage...${NC}"
if [ -f "$TEST_DATA_DIR/weight.json" ]; then
    echo -e "${GREEN}Data file created at: $TEST_DATA_DIR/weight.json${NC}"
    echo "Content:"
    cat "$TEST_DATA_DIR/weight.json"
else
    echo -e "${RED}Data file not found!${NC}"
fi