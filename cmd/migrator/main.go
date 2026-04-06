package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

// helper to safely read env variables
func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env variable: %s", key)
	}
	return v
}

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// command flags
	action := flag.String("action", "up", "up | down | force | version")
	steps := flag.Int("steps", 1, "number of steps for down/up when applicable")
	forceVersion := flag.Int("forceVersion", -1, "force database version (used when action=force)")
	flag.Parse()

	// read database config from .env
	dbUser := getEnv("DB_USER")
	dbPass := getEnv("DB_PASSWORD")
	dbHost := getEnv("DB_HOST")
	dbPort := getEnv("DB_PORT")
	dbName := getEnv("DB_NAME")
	dbSSL := getEnv("DB_SSLMODE")

	// build database URL
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
		dbSSL,
	)

	// migration folder path
	// migrationsPath := "file://E:/QStack-Backend/migrations"
	migrationsPath := "file://./migrations"

	// create migrator
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		log.Fatal("migrate.New error: ", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	// actions
	switch *action {
	case "up":
		err = m.Up()

	case "down":
		err = m.Steps(-*steps)

	case "version":
		v, dirty, verr := m.Version()
		if verr != nil {
			log.Fatal("version error: ", verr)
		}
		fmt.Printf("version=%d dirty=%v\n", v, dirty)
		return

	case "force":
		if *forceVersion < 0 {
			log.Fatal("forceVersion must be >= 0")
		}
		err = m.Force(*forceVersion)

	default:
		log.Fatal("unknown action: ", *action)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration error: ", err)
	}

	log.Println("migration done:", *action)
}
