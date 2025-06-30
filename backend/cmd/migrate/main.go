package main

import (
	"fmt"
	"log"
	"os"

	"reminder-app/db"
	"reminder-app/db/migrate"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/reminder?sslmode=disable"
	}

	// Connect to database
	dbConnections, err := db.New(dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConnections.Close()

	// Create migrator
	migrator := migrate.NewMigrator(dbConnections.GORM)

	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up", "migrate":
		if err := migrator.Migrate(); err != nil {
			log.Fatal("Migration failed:", err)
		}
		fmt.Println("All migrations completed successfully!")

	case "rollback":
		migrationID := ""
		if len(os.Args) >= 3 {
			migrationID = os.Args[2]
		}
		if err := migrator.RollbackTo(migrationID); err != nil {
			log.Fatal("Rollback failed:", err)
		}
		fmt.Println("Rollback completed successfully!")

	case "rollback-last":
		if err := migrator.RollbackLast(); err != nil {
			log.Fatal("Rollback failed:", err)
		}
		fmt.Println("Last migration rolled back successfully!")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go up              - Run all pending migrations")
	fmt.Println("  go run cmd/migrate/main.go migrate         - Alias for 'up'")
	fmt.Println("  go run cmd/migrate/main.go rollback [ID]   - Rollback to migration ID")
	fmt.Println("  go run cmd/migrate/main.go rollback-last   - Rollback last migration")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go up")
	fmt.Println("  go run cmd/migrate/main.go rollback 202412291545")
	fmt.Println("  go run cmd/migrate/main.go rollback-last")
}