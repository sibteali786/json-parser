#!/bin/bash
echo "Testing all fail*.json files..."
echo "================================"

passed=0
failed=0

for file in tests/step5/fail*.json; do
    filename=$(basename "$file")
    
    result=$(./jsonparser "$file" 2>&1)
    exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        echo "❌ $filename: UNEXPECTED PASS"
        echo "   Content: $(head -1 "$file" | tr -d '\n')"
        echo "   Parser output: $(echo "$result" | head -1)"
        ((failed++))
    else
        echo "✅ $filename: CORRECTLY FAILED"
        echo "   Error: $(echo "$result" | grep "Error:" | head -1)"
        ((passed++))
    fi
    echo
done

echo "================================"
echo "Results: $passed correct, $failed unexpected passes"