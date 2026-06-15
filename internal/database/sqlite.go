package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

func Open(dsn string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dsn), 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	if err := seed(db); err != nil {
		return nil, fmt.Errorf("seed: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	PRAGMA journal_mode=WAL;
	PRAGMA foreign_keys=ON;

	CREATE TABLE IF NOT EXISTS users (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		username  TEXT    NOT NULL UNIQUE,
		password  TEXT    NOT NULL,
		role      TEXT    NOT NULL DEFAULT 'master',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS orders (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		client_name     TEXT    NOT NULL,
		phone           TEXT    NOT NULL,
		device          TEXT    NOT NULL,
		description     TEXT    NOT NULL,
		estimated_cost  REAL    NOT NULL DEFAULT 0,
		status          TEXT    NOT NULL DEFAULT 'accepted',
		created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS parts (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		name          TEXT    NOT NULL,
		quantity      INTEGER NOT NULL DEFAULT 0,
		purchase_price REAL   NOT NULL DEFAULT 0,
		created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS order_parts (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id  INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		part_id   INTEGER NOT NULL REFERENCES parts(id),
		quantity  INTEGER NOT NULL DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		token  TEXT PRIMARY KEY,
		data   BLOB NOT NULL,
		expiry REAL NOT NULL
	);

	CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry);
	`

	_, err := db.Exec(schema)
	return err
}

func seed(db *sql.DB) error {
	users := []struct {
		username string
		password string
		role     string
	}{
		{"admin", "admin123", "admin"},
		{"master", "master123", "master"},
	}

	for _, u := range users {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", u.username).Scan(&count); err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		if _, err := db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", u.username, string(hash), u.role); err != nil {
			return err
		}
	}

	return nil
}
