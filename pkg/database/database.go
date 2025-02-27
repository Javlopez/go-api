package database

import (
	"fmt"
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

func (db *Database) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.config.Host, db.config.Port, db.config.User, db.config.Password, db.config.DBName, db.config.SSLMode,
	)
}

func (db *Database) Connect() (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", db.DSN())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
