# scripts/run_tests.sh
#!/bin/bash

# Make all test scripts executable
chmod +x scripts/test_*.sh

# Initialize counters
total=0
passed=0
failed=0

# Run all test scripts
for test in scripts/test_*.sh; do
    if [[ "$test" != *"framework.sh" ]]; then
        echo "Running $test..."
        if $test; then
            ((passed++))
        else
            ((failed++))
        fi
        ((total++))
    fi
done

# Show summary
echo "Test Summary:"
echo "Total: $total"
echo "Passed: $passed"
echo "Failed: $failed"

exit $failed