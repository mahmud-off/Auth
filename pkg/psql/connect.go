package psql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresConnectionConfig struct {
	Name     string
	Password string
	Address  string
	Port     string
	SSLmode  string
}

func NewPostgresDB(config *PostgresConnectionConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s?sslmode=%s", config.Name, config.Password, config.Address, config.Port, config.SSLmode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = PingPostgres(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func PingPostgres(db *sqlx.DB) error {
	err := db.Ping()
	return err
}
