package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rcovery/go-stock-control/internal/config"
	"github.com/rcovery/go-stock-control/internal/http/handlers"
	database_turso "github.com/rcovery/go-stock-control/internal/infra/database/turso"
	"github.com/rcovery/go-stock-control/internal/part"
	part_repository "github.com/rcovery/go-stock-control/internal/part/repository/turso"
)

func main() {
	config.InitConfig()

	connectionString := database_turso.GetConnectionFromEnv()
	db, databaseErr := database_turso.NewDatabaseConnection(connectionString)
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

	host := config.GetString("HOST")
	port := config.GetString("PORT")

	log.Printf("starting server on %s:%s", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil))
}
