# scripts/test_weight_delete.sh
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

TEST_DATA_DIR=~/.health-tracker/data/test

# Cleanup and setup
rm -rf "$TEST_DATA_DIR"

echo -e "${YELLOW}Starting weight delete tests...${NC}"

# Setup test data
echo -e "${YELLOW}Setting up test data...${NC}"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "First weight"
echo "y" | TEST_MODE=true ./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "Second weight"

# Test successful deletion
echo -e "${YELLOW}Test 1: Successful deletion${NC}"
echo -e "${YELLOW}Before delete - current records:${NC}"
TEST_MODE=true ./bin/tracker weight list --from 2024-01-01 --to 2024-01-31

echo "y" | TEST_MODE=true ./bin/tracker weight delete w00001

echo -e "\n${YELLOW}Verifying deletion:${NC}"
remaining=$(TEST_MODE=true ./bin/tracker weight list --from 2024-01-01 --to 2024-01-31)
echo "$remaining"
if [[ "$remaining" != *"w00001"* ]] && [[ "$remaining" == *"w00002"* ]]; then
    echo -e "${GREEN}✓ Record successfully deleted${NC}"
else
    echo -e "${RED}✗ Delete verification failed${NC}"
fi

# Test non-existent record
echo -e "\n${YELLOW}Test 2: Non-existent record${NC}"
output=$(TEST_MODE=true ./bin/tracker weight delete w99999 2>&1)
if [[ "$output" == *"not found"* ]]; then
    echo -e "${GREEN}✓ Correctly handled non-existent record${NC}"
    echo -e "Output: $output"
else
    echo -e "${RED}✗ Failed to handle non-existent record${NC}"
    echo -e "Output: $output"
fi

# Test delete cancellation
echo -e "\n${YELLOW}Test 3: Testing delete cancellation${NC}"
output=$(echo "n" | TEST_MODE=true ./bin/tracker weight delete w00002 2>&1)
if [[ "$output" == *"cancelled"* ]]; then
    echo -e "${GREEN}✓ Correctly handled deletion cancellation${NC}"
    # Verify record still exists
    remaining=$(TEST_MODE=true ./bin/tracker weight list --from 2024-01-01 --to 2024-01-31)
    if [[ "$remaining" == *"w00002"* ]]; then
        echo -e "${GREEN}✓ Record still exists after cancellation${NC}"
    else
        echo -e "${RED}✗ Record missing after cancellation${NC}"
    fi
else
    echo -e "${RED}✗ Failed to handle deletion cancellation${NC}"
    echo -e "Output: $output"
fi

# Cleanup
rm -rf "$TEST_DATA_DIR"