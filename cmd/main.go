package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rcovery/go-stock-control/internal/config"
	"github.com/rcovery/go-stock-control/internal/http/handlers"
	database_postgres "github.com/rcovery/go-stock-control/internal/infra/database/postgres"
	"github.com/rcovery/go-stock-control/internal/part"
	"github.com/rcovery/go-stock-control/internal/part/postgres"
)

func main() {
	config.InitConfig()

	connectionString := database_postgres.GetConnectionFromEnv()
	db, databaseErr := database_postgres.NewDatabaseConnection(connectionString)
	if databaseErr != nil {
		panic(databaseErr)
	}

	repoInstance := postgres.NewRepository(db)
	serviceInstance := part.NewService(repoInstance)
	handlers.HandlePart(serviceInstance)

	host := config.GetString("HOST")
	port := config.GetString("PORT")

	log.Printf("starting server on %s:%s", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil))
}
