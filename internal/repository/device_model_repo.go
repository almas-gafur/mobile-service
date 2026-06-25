package repository

import (
	"database/sql"
	"mobile-service/internal/models"
)

type DeviceModelRepo struct {
	db *sql.DB
}

func NewDeviceModelRepo(db *sql.DB) *DeviceModelRepo {
	return &DeviceModelRepo{db: db}
}

func (r *DeviceModelRepo) List() ([]models.DeviceModel, error) {
	rows, err := r.db.Query("SELECT id, name, created_at FROM device_models ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mods []models.DeviceModel
	for rows.Next() {
		var m models.DeviceModel
		if err := rows.Scan(&m.ID, &m.Name, &m.CreatedAt); err != nil {
			return nil, err
		}
		mods = append(mods, m)
	}
	return mods, rows.Err()
}

func (r *DeviceModelRepo) GetByID(id int64) (*models.DeviceModel, error) {
	m := &models.DeviceModel{}
	err := r.db.QueryRow("SELECT id, name, created_at FROM device_models WHERE id = ?", id).Scan(&m.ID, &m.Name, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

func (r *DeviceModelRepo) Create(m *models.DeviceModel) (int64, error) {
	res, err := r.db.Exec("INSERT INTO device_models (name) VALUES (?)", m.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *DeviceModelRepo) Update(m *models.DeviceModel) error {
	_, err := r.db.Exec("UPDATE device_models SET name = ? WHERE id = ?", m.Name, m.ID)
	return err
}

func (r *DeviceModelRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM device_models WHERE id = ?", id)
	return err
}
