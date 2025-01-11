# scripts/test_data.sh
#!/bin/bash

# Function to load a standard set of test data
load_test_data() {
    ./bin/tracker weight add -v 185.5 --date 2024-01-08 --notes "Morning weight"
    ./bin/tracker weight add -v 184.8 --date 2024-01-09 --notes "After workout"
    ./bin/tracker weight add -v 184.2 --date 2024-01-10 --notes "Progress"
}

# Function to load specific test scenarios
load_scenario_weight_change() {
    ./bin/tracker weight add -v 185.5 --date 2024-01-08
    ./bin/tracker weight add -v 175.5 --date 2024-01-09  # Should trigger warning
}