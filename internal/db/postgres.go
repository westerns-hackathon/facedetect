package db

import (
	"face-detection/internal/config"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импортируем PostgreSQL драйвер
)

type Database struct {
	db  *sqlx.DB
	cfg *config.Config
}

func NewDB(cfg *config.Config) *Database {
	return &Database{cfg: cfg}
}

func (d *Database) Open() error {
	dsn := fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable password=%s host=%s port=5432",
		d.cfg.Database.Username, d.cfg.Database.DBName, d.cfg.Database.Password,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	d.db = db
	log.Println("Successfully connected to postgres database")
	return nil
}

func (d *Database) Close() error {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	return nil
}
