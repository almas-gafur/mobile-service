package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/example/repair-crm/internal/models"
)

type MasterRepository struct {
	db *sql.DB
}

func NewMasterRepository(db *sql.DB) *MasterRepository {
	return &MasterRepository{db: db}
}

func (r *MasterRepository) FindByUsername(ctx context.Context, username string) (*models.Master, error) {
	const query = `
		SELECT m.id, m.workshop_id, w.name, m.username, m.password_hash
		FROM masters m
		JOIN workshops w ON w.id = m.workshop_id
		WHERE m.username = $1
	`

	var master models.Master
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&master.ID,
		&master.WorkshopID,
		&master.WorkshopName,
		&master.Username,
		&master.PasswordHash,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &master, nil
}
