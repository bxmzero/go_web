package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"go.uber.org/fx"

	"example.com/gin-fx-sqlite-demo/internal/config"
)

func NewSQLite(lifecycle fx.Lifecycle, cfg config.Config) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := initializeSchema(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := seedData(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

func initializeSchema(db *sql.DB) error {
	const schema = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    item TEXT NOT NULL,
    amount INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("initialize schema: %w", err)
	}

	return nil
}

func seedData(db *sql.DB) error {
	var userCount int
	if err := db.QueryRow(`SELECT COUNT(1) FROM users`).Scan(&userCount); err != nil {
		return fmt.Errorf("count users: %w", err)
	}

	if userCount == 0 {
		if _, err := db.Exec(`
INSERT INTO users (name, email) VALUES
    ('Alice', 'alice@example.com'),
    ('Bob', 'bob@example.com');`); err != nil {
			return fmt.Errorf("seed users: %w", err)
		}
	}

	var orderCount int
	if err := db.QueryRow(`SELECT COUNT(1) FROM orders`).Scan(&orderCount); err != nil {
		return fmt.Errorf("count orders: %w", err)
	}

	if orderCount == 0 {
		if _, err := db.Exec(`
INSERT INTO orders (user_id, item, amount) VALUES
    (1, 'Keyboard', 299),
    (1, 'Mouse', 129),
    (2, 'Monitor', 1499);`); err != nil {
			return fmt.Errorf("seed orders: %w", err)
		}
	}

	return nil
}
