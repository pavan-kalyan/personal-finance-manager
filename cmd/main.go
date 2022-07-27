package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/pkger"
	"github.com/markbates/pkger"
	"go-personal-finance/pkg/api"
	"go-personal-finance/pkg/app"
	"go-personal-finance/pkg/repository"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to start app. error: %s\n", err)
	}
}

func run() error {

	shouldMigratePtr := flag.Bool("migrate", false, "apply any new migration files. Only supports forward migrations.")
	dbFilePathPtr := flag.String("db-path", "./sqlite.db", "relative file path from current folder.")

	flag.Parse()

	db, err := setupDatabase(shouldMigratePtr, dbFilePathPtr)
	if err != nil {
		return err
	}

	// setup router dependency
	router := gin.Default()
	//TODO: this is not secure, remove it later
	router.Use(cors.Default())

	accountsDAO := repository.NewAccountsDAO(db)
	// create all required services
	accountsService := api.NewAccountsService(accountsDAO)

	server := app.NewServer(router, accountsService)
	err = server.Run()

	return err
}

func setupDatabase(shouldMigratePtr *bool, dbFilePathPtr *string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", *dbFilePathPtr)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database with err: %w\n", err)
		//log.Fatalf("Failed to open database with err: %v\n", err)
	}

	// ensuring that the db is connected.
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping database with err: %w\n", err)
		//log.Fatalf("Failed to ping database with err: %v\n", err)
	}

	err = runMigrations(shouldMigratePtr, dbFilePathPtr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func runMigrations(shouldMigratePtr *bool, dbFilePathPtr *string) error {
	if *shouldMigratePtr {
		_ = pkger.Include("../migrations")
		migrator, err := migrate.New("pkger://../migrations", fmt.Sprintf("sqlite://%s", *dbFilePathPtr))
		if err != nil {
			return fmt.Errorf("Could not access sqlite db with err: %v\n", err)
		}

		err = migrator.Up()
		if err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("Failed to run migrations with err: %v\n", err)
		}
	}
	return nil
}
