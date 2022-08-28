package postgres

import (
	"context"
	"fmt"
	"txsystem/config"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	// Подключаемся к сиситемной роли в постгрес
	databaseUrl := fmt.Sprintf(`postgres://%s:%s@%s:%d/%s`,
		cfg.Username,
		cfg.Password,
		cfg.Hostname,
		cfg.Port,
		"postgres")

	pool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	// Проверяем на существование БД
	var exists bool
	if err := pool.QueryRow(context.Background(), `SELECT COUNT(*)>0 AS db_exists FROM pg_database WHERE datname = $1;`, cfg.Database).Scan(&exists); err != nil {
		return nil, fmt.Errorf("cannot check if database exists: %v", err)
	}

	if !exists {
		if err := createDB(pool, cfg.Database); err != nil {
			return nil, fmt.Errorf("cannot init database: %v", err)
		}
	}
	pool.Close()

	// Инициализируем подключение к БД
	databaseUrl = fmt.Sprintf(`postgres://%s:%s@%s:%d/%s`,
		cfg.Username,
		cfg.Password,
		cfg.Hostname,
		cfg.Port,
		cfg.Database)

	pool, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	if err := createTables(pool); err != nil {
		return nil, err
	}

	return pool, nil
}

func createDB(pool *pgxpool.Pool, dbName string) error {
	query := `CREATE DATABASE ` + dbName + ";"
	_, err := pool.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}

func createTables(pool *pgxpool.Pool) error {
	users := `--sql
		CREATE TABLE IF NOT EXISTS users (
		"id" INT NOT NULL,
		"balance" INT,
		PRIMARY KEY ("id")
	);`

	txs := `--sql
		CREATE TABLE IF NOT EXISTS txs (
		"id" SERIAL,
		"user_id" INT NOT NULL,
		"amount" INT NOT NULL,
		"action" TEXT NOT NULL,
		"create_at" TIMESTAMP DEFAULT now(),
		"status" INT NOT NULL,
		PRIMARY KEY ("id"),
		CONSTRAINT "fk_user_id"
			FOREIGN KEY ("user_id")
				REFERENCES users("id")
	);`

	_, err := pool.Exec(context.Background(), users+txs)
	if err != nil {
		return errors.Wrap(err, "create table error")
	}

	return nil
}
