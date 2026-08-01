package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rcovery/go-stock-control/internal/config"
	"github.com/rcovery/go-stock-control/internal/http/handlers"
	"github.com/rcovery/go-stock-control/internal/infra/database/postgres"
	"github.com/rcovery/go-stock-control/internal/part"
	part_repository "github.com/rcovery/go-stock-control/internal/part/repository/postgres"
)

func main() {
	config.InitConfig()

	connectionString := postgres.GetConnectionFromEnv()
	db, databaseErr := postgres.NewDatabaseConnection(connectionString)
	if databaseErr != nil {
		panic(databaseErr)
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
