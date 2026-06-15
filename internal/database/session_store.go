package database

import (
	"context"
	"database/sql"
	"time"
)

// SQLiteSessionStore implements scs.Store backed by SQLite.
type SQLiteSessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SQLiteSessionStore {
	return &SQLiteSessionStore{db: db}
}

func (s *SQLiteSessionStore) Find(token string) ([]byte, bool, error) {
	var data []byte
	var expiry float64
	err := s.db.QueryRow(
		"SELECT data, expiry FROM sessions WHERE token = ?", token,
	).Scan(&data, &expiry)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if float64(time.Now().Unix()) > expiry {
		return nil, false, nil
	}
	return data, true, nil
}

func (s *SQLiteSessionStore) Commit(token string, data []byte, expiry time.Time) error {
	_, err := s.db.Exec(
		`INSERT INTO sessions (token, data, expiry) VALUES (?, ?, ?)
		 ON CONFLICT(token) DO UPDATE SET data=excluded.data, expiry=excluded.expiry`,
		token, data, float64(expiry.Unix()),
	)
	return err
}

func (s *SQLiteSessionStore) Delete(token string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

// CleanupExpired removes stale sessions. Call this periodically.
func (s *SQLiteSessionStore) CleanupExpired(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.db.Exec("DELETE FROM sessions WHERE expiry < ?", float64(time.Now().Unix()))
		}
	}
}
