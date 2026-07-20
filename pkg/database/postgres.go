package database

import (
	"database/sql"
	"fmt"
	"log"
	"task-manager/pkg/config"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var db *sql.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		log.Println("try to access database failed, retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL,
		assignee VARCHAR(255),
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}
