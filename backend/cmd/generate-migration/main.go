package main

import (
	_ "embed"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"reminder-app/config"

	"go.uber.org/fx"
)

//go:embed new_migration.tmpl
var migrationTemplate string

type MigrationData struct {
	Timestamp       string
	Description     string
	StructName      string
	PlanName        string
	ConstructorName string
	UpFuncName      string
	DownFuncName    string
}

type MigrationGenerator struct {
	config *config.Config
}

func NewMigrationGenerator(cfg *config.Config) *MigrationGenerator {
	return &MigrationGenerator{config: cfg}
}

func (mg *MigrationGenerator) Generate(description string) error {
	// Generate timestamp
	timestamp := time.Now().Format("200601021504") // YYYYMMDDHHMM

	// Create migration data
	data := MigrationData{
		Timestamp:       timestamp,
		Description:     description,
		StructName:      fmt.Sprintf("Migration%s%s", timestamp, description),
		PlanName:        fmt.Sprintf("MigrationPlan%s", timestamp),
		ConstructorName: fmt.Sprintf("NewMigrationPlan%s", timestamp),
		UpFuncName:      fmt.Sprintf("Up%s", timestamp),
		DownFuncName:    fmt.Sprintf("Down%s", timestamp),
	}

	// Parse embedded template
	tmpl, err := template.New("migration").Parse(migrationTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	outputFilename := fmt.Sprintf("%s_%s.go", timestamp, description)
	outputPath := filepath.Join("db", "migrate", outputFilename)

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("✅ Created migration file: %s\n", outputPath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("1. Edit %s and implement Up%s() and Down%s() methods\n", outputPath, timestamp, timestamp)
	fmt.Printf("2. Add Plan%s to the plans slice in db/migrate/migrator.go\n", timestamp)
	fmt.Println("3. Run migrations with: go run ./cmd/migrate up")

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/generate-migration/main.go <description>")
		fmt.Println("Example: go run cmd/generate-migration/main.go AddUserPreferences")
		os.Exit(1)
	}

	description := os.Args[1]

	// Create fx app for dependency injection
	app := fx.New(
		config.Module,
		fx.Provide(NewMigrationGenerator),
		fx.Invoke(func(generator *MigrationGenerator, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return generator.Generate(description)
				},
			})
		}),
		fx.NopLogger, // Suppress fx logs for cleaner output
	)

	// Start and stop the app immediately
	if err := app.Start(context.Background()); err != nil {
		log.Fatal("Failed to generate migration:", err)
	}
	if err := app.Stop(context.Background()); err != nil {
		log.Fatal("Failed to stop app:", err)
	}
}
