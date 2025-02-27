package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	config *Config
}

func New(config *Config) *Database {
	return &Database{
		config: config,
	}
}

func (db *Database) Connect() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", db.config.DSN())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
