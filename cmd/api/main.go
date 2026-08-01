package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rcovery/go-stock-control/internal/http/handlers"
	database_turso "github.com/rcovery/go-stock-control/internal/infra/database/turso"
	"github.com/rcovery/go-stock-control/internal/part"
	part_repository "github.com/rcovery/go-stock-control/internal/part/repository/turso"
)

const (
	host = "0.0.0.0"
	port = "9000"
)

func main() {
	db, databaseErr := database_turso.NewDatabaseConnection(database_turso.DatabasePath)
	if databaseErr != nil {
		panic(databaseErr)
	}

	if migrateErr := database_turso.Migrate(context.Background(), db); migrateErr != nil {
		panic(migrateErr)
	}

	repoInstance := part_repository.NewRepository(db)
	serviceInstance := part.NewService(repoInstance)
	partHandler := handlers.NewPartHandler(serviceInstance)
	partHandler.HandlePart()

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
