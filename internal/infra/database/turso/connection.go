package turso

import (
	"database/sql"

	_ "turso.tech/database/tursogo"
)

const DatabasePath = "app.db"

func NewDatabaseConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("turso", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
