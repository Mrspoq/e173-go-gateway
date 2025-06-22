#!/bin/bash

# Script to convert enterprise repositories from *sqlx.DB to *pgxpool.Pool

echo "Converting enterprise repositories to use *pgxpool.Pool..."

REPO_FILES=(
    "/root/e173_go_gateway/internal/repository/customer_repository.go"
    "/root/e173_go_gateway/internal/repository/payment_repository.go"
    "/root/e173_go_gateway/internal/repository/routing_repository.go"
    "/root/e173_go_gateway/internal/repository/system_repository.go"
)

for file in "${REPO_FILES[@]}"; do
    echo "Processing $file..."
    
    # Create backup
    cp "$file" "$file.backup"
    
    # Update imports
    sed -i 's|"database/sql"|"context"|g' "$file"
    sed -i 's|"github.com/jmoiron/sqlx"|"github.com/jackc/pgx/v5"\n\t"github.com/jackc/pgx/v5/pgxpool"|g' "$file"
    
    # Update struct field
    sed -i 's|db \*sqlx\.DB|db *pgxpool.Pool|g' "$file"
    
    # Update constructor
    sed -i 's|db \*sqlx\.DB|db *pgxpool.Pool|g' "$file"
    
    # Update database method calls
    sed -i 's|sql\.ErrNoRows|pgx.ErrNoRows|g' "$file"
    sed -i 's|r\.db\.NamedQuery(|rows, err := r.db.Query(context.Background(), |g' "$file"
    sed -i 's|r\.db\.NamedExec(|_, err := r.db.Exec(context.Background(), |g' "$file"
    sed -i 's|r\.db\.Get(|err := r.db.QueryRow(context.Background(), |g' "$file"
    sed -i 's|r\.db\.Select(|rows, err := r.db.Query(context.Background(), |g' "$file"
    sed -i 's|r\.db\.Exec(|_, err := r.db.Exec(context.Background(), |g' "$file"
    
    echo "Converted $file"
done

echo "Repository conversion completed!"
echo "Note: Manual review and fixes may be needed for complex queries."
