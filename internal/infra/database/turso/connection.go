package turso

import (
	"database/sql"

	"github.com/rcovery/go-stock-control/internal/config"
	_ "turso.tech/database/tursogo"
)

func GetConnectionFromEnv() string {
	path := config.GetString("TURSO_DATABASE_PATH")
	if path == "" {
		path = "app.db"
	}

	return path
}

func NewDatabaseConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("turso", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
