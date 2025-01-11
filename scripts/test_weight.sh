# 1. First, let's remove the existing data
rm -rf ~/.health-tracker/data/production/

# 2. Build the latest version
go build -o bin/tracker ./cmd/tracker

# 3. Test the help commands
./bin/tracker -h
./bin/tracker weight -h
./bin/tracker weight add -h

# 3. Let's run a series of test commands:
# First record
./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "Morning weight"

# Verify the data was stored
cat ~/.health-tracker/data/production/weight.json

# Add another record for the next day
./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "Good progress"

# Check the data file again to verify IDs and data
cat ~/.health-tracker/data/production/weight.json

# 4.validation tests 
# Invalid weight (too low)
./bin/tracker weight add -v 45.0 --date 2024-01-11

# Invalid date format
./bin/tracker weight add -v 185.0 --date "Jan 11, 2024"

# Future date
./bin/tracker weight add -v 185.0 --date 2025-01-01

# Test data validation with invalid weight
./bin/tracker weight add -v 45.0 --date 2024-01-10

# - Add a duplicate date
./bin/tracker weight add -v 184.0 --date 2024-01-09 --notes "This should prompt for overwrite"

# - Add a significant weight change
./bin/tracker weight add -v 175.0 --date 2024-01-10 --notes "This should warn about large change"

# This should trigger a warning about significant weight change
./bin/tracker weight add -v 2255.0 --date 2024-01-10 --notes "Big change - should warn"

# 5. Test listing records
./bin/tracker weight list
./bin/tracker weight list --from 2024-01-01 --to 2024-01-10

# Test getting a specific record
./bin/tracker weight get --date 2024-01-08

# Test updating a record (use an ID from your list output)
./bin/tracker weight update w00001 --value 185.0

# Test deleting a record (use an ID from your list output)
./bin/tracker weight delete w00001
