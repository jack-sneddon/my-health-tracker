# scripts/test_weight_update.sh
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

TEST_DATA_DIR=~/.health-tracker/data/test

# Cleanup and setup
rm -rf "$TEST_DATA_DIR"

echo -e "${YELLOW}Starting weight update tests...${NC}"

# Setup test data
echo -e "${YELLOW}Setting up test data...${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "Second weight"

# Show initial state
echo -e "${YELLOW}Before update - current records:${NC}"
TEST_MODE=true ./bin/tracker weight list --from 2024-01-01 --to 2024-01-31

# Test 1: Update weight value
echo -e "\n${YELLOW}Test 1: Update weight value${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight update w00001 --value 186.0

# Verify update
echo -e "\n${YELLOW}Verifying weight update:${NC}"
output=$(TEST_MODE=true ./bin/tracker weight list --from 2024-01-01 --to 2024-01-31)
echo "$output"
if [[ "$output" == *"186.0"* ]]; then
    echo -e "${GREEN}✓ Weight value updated successfully${NC}"
else
    echo -e "${RED}✗ Weight update failed${NC}"
fi

# Test 2: Update notes
echo -e "\n${YELLOW}Test 2: Update notes${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight update w00001 --notes "Updated notes"

# Verify notes update
echo -e "\n${YELLOW}Verifying notes update:${NC}"
output=$(TEST_MODE=true ./bin/tracker weight list --from 2024-01-01 --to 2024-01-31)
echo "$output"
if [[ "$output" == *"Updated notes"* ]]; then
    echo -e "${GREEN}✓ Notes updated successfully${NC}"
else
    echo -e "${RED}✗ Notes update failed${NC}"
fi

# Test 3: Invalid weight value
echo -e "\n${YELLOW}Test 3: Invalid weight value${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w00001 --value 45.0 2>&1)
if [[ "$output" == *"validation failed"* ]]; then
    echo -e "${GREEN}✓ Correctly handled invalid weight value${NC}"
else
    echo -e "${RED}✗ Failed to handle invalid weight value${NC}"
    echo "Output: $output"
fi

# Test 4: Non-existent record
echo -e "\n${YELLOW}Test 4: Non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker weight update w99999 --value 185.0 2>&1)
if [[ "$output" == *"not found"* ]]; then
    echo -e "${GREEN}✓ Correctly handled non-existent record${NC}"
else
    echo -e "${RED}✗ Failed to handle non-existent record${NC}"
    echo "Output: $output"
fi

# Cleanup
rm -rf "$TEST_DATA_DIR"