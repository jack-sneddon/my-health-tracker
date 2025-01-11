# scripts/cleanup.sh
#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Starting cleanup...${NC}"

# List of files to remove
files=(
    "internal/config/config.go"
    "internal/reporter/reporter.go"
    "pkg/utils/utils.go"
    "web/templates/dashboard.html"
    "scripts/test_weight_orig.sh"
    "scripts/test_weight_crud.sh"
    "tracker"
)

# List of directories to remove if empty
directories=(
    "internal/config"
    "internal/reporter"
    "pkg/utils"
    "web/templates"
    "web"
)

# Remove files
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        rm "$file"
        echo -e "${GREEN}Removed file: $file${NC}"
    else
        echo -e "${RED}File not found: $file${NC}"
    fi
done

# Remove empty directories
for dir in "${directories[@]}"; do
    if [ -d "$dir" ] && [ -z "$(ls -A $dir)" ]; then
        rmdir "$dir"
        echo -e "${GREEN}Removed empty directory: $dir${NC}"
    elif [ -d "$dir" ]; then
        echo -e "${YELLOW}Directory not empty, skipping: $dir${NC}"
    else
        echo -e "${RED}Directory not found: $dir${NC}"
    fi
done

echo -e "${GREEN}Cleanup complete${NC}"