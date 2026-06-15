package repository

import (
	"database/sql"
	"fmt"

	"repair-crm/internal/models"
)

type PartRepo struct {
	db *sql.DB
}

func NewPartRepo(db *sql.DB) *PartRepo {
	return &PartRepo{db: db}
}

func (r *PartRepo) List() ([]models.Part, error) {
	rows, err := r.db.Query(
		"SELECT id, name, quantity, purchase_price, created_at FROM parts ORDER BY name ASC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []models.Part
	for rows.Next() {
		var p models.Part
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.PurchasePrice, &p.CreatedAt); err != nil {
			return nil, err
		}
		parts = append(parts, p)
	}
	return parts, rows.Err()
}

func (r *PartRepo) GetByID(id int64) (*models.Part, error) {
	p := &models.Part{}
	err := r.db.QueryRow(
		"SELECT id, name, quantity, purchase_price, created_at FROM parts WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Name, &p.Quantity, &p.PurchasePrice, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, err
}

func (r *PartRepo) Create(p *models.Part) (int64, error) {
	res, err := r.db.Exec(
		"INSERT INTO parts (name, quantity, purchase_price) VALUES (?, ?, ?)",
		p.Name, p.Quantity, p.PurchasePrice,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PartRepo) Update(p *models.Part) error {
	_, err := r.db.Exec(
		"UPDATE parts SET name = ?, quantity = ?, purchase_price = ? WHERE id = ?",
		p.Name, p.Quantity, p.PurchasePrice, p.ID,
	)
	return err
}

// WriteOff deducts qty from part and links it to an order atomically.
func (r *PartRepo) WriteOff(orderID, partID int64, qty int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var available int
	if err := tx.QueryRow("SELECT quantity FROM parts WHERE id = ?", partID).Scan(&available); err != nil {
		return err
	}
	if available < qty {
		return fmt.Errorf("недостаточно запчастей на складе: доступно %d, запрошено %d", available, qty)
	}

	if _, err := tx.Exec("UPDATE parts SET quantity = quantity - ? WHERE id = ?", qty, partID); err != nil {
		return err
	}

	if _, err := tx.Exec(
		"INSERT INTO order_parts (order_id, part_id, quantity) VALUES (?, ?, ?)",
		orderID, partID, qty,
	); err != nil {
		return err
	}

	return tx.Commit()
}
