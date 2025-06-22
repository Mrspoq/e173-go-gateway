#!/bin/bash

# Script to fix duplicate variable declarations caused by automated conversion

echo "Fixing syntax errors in converted repositories..."

REPO_FILES=(
    "/root/e173_go_gateway/internal/repository/customer_repository.go"
    "/root/e173_go_gateway/internal/repository/payment_repository.go"
    "/root/e173_go_gateway/internal/repository/routing_repository.go"
    "/root/e173_go_gateway/internal/repository/system_repository.go"
)

for file in "${REPO_FILES[@]}"; do
    echo "Fixing syntax errors in $file..."
    
    # Fix duplicate variable declarations
    sed -i 's|rows, err := rows, err := |rows, err := |g' "$file"
    sed -i 's|_, err := _, err := |_, err := |g' "$file"
    sed -i 's|err := err := |err := |g' "$file"
    
    # Fix any remaining duplicate patterns
    sed -i 's|rows, err := r\.db\.Query(context\.Background(), rows, err := r\.db\.Query(context\.Background(),|rows, err := r.db.Query(context.Background(),|g' "$file"
    sed -i 's|_, err := r\.db\.Exec(context\.Background(), _, err := r\.db\.Exec(context\.Background(),|_, err := r.db.Exec(context.Background(),|g' "$file"
    
    echo "Fixed $file"
done

echo "Syntax error fixes completed!"
