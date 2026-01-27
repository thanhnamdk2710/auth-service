package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/thanhnamdk2710/auth-service/internal/config"
)

func main() {
	var direction string
	var steps int

	flag.StringVar(&direction, "direction", "up", "Migration direction: up, down, version, force")
	flag.IntVar(&steps, "steps", 0, "Number of migrations to apply (0 = all)")
	flag.Parse()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up completed successfully")

	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down completed successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		log.Printf("Current version: %d, Dirty: %t", version, dirty)

	case "force":
		if steps == 0 {
			log.Fatal("Force requires a version number via -steps flag")
		}
		err = m.Force(steps)
		if err != nil {
			log.Fatalf("Force failed: %v", err)
		}
		log.Printf("Forced to version %d", steps)

	default:
		log.Fatalf("Unknown direction: %s", direction)
	}
}
