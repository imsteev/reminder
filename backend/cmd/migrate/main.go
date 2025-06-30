package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"reminder-app/config"
	gormmodule "reminder-app/db/gorm"
	"reminder-app/db/migrate"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Create fx app for dependency injection
	app := fx.New(
		config.Module,
		gormmodule.Module,
		fx.Invoke(func(gormDB *gorm.DB, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					// Create migrator
					migrator := migrate.NewMigrator(gormDB)

					// Execute migration command
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

					return nil
				},
				OnStop: func(ctx context.Context) error {
					// Close database connections
					if sqlDB, err := gormDB.DB(); err == nil {
						sqlDB.Close()
					}
					return nil
				},
			})
		}),
		fx.NopLogger, // Suppress fx logs for cleaner output
	)

	// Start and stop the app immediately
	if err := app.Start(context.Background()); err != nil {
		log.Fatal("Failed to start app:", err)
	}
	if err := app.Stop(context.Background()); err != nil {
		log.Fatal("Failed to stop app:", err)
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