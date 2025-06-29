package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://e173_user:3omartel580@localhost/e173_gateway?sslmode=disable"
	}

	// Connect to database
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Check if admin user already exists
	var exists bool
	err = db.QueryRow(context.Background(), 
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", "admin").Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check if admin exists: %v", err)
	}

	if exists {
		fmt.Println("Admin user already exists")
		
		// Update the password
		_, err = db.Exec(context.Background(),
			"UPDATE users SET password_hash = $1 WHERE username = $2",
			string(passwordHash), "admin")
		if err != nil {
			log.Fatalf("Failed to update admin password: %v", err)
		}
		fmt.Println("Admin password updated successfully")
	} else {
		// Create admin user
		_, err = db.Exec(context.Background(), `
			INSERT INTO users (username, email, password_hash, first_name, last_name, role, is_active, is_2fa_enabled)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			"admin", "admin@e173gateway.local", string(passwordHash), "System", "Administrator", "super_admin", true, false)
		
		if err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}
		fmt.Println("Admin user created successfully")
	}

	fmt.Println("Username: admin")
	fmt.Println("Password: admin")
}