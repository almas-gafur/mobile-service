package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"mobile-service/internal/"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) List(status, search string) ([]models.Order, error) {
	query := "SELECT id, client_name, phone, device, description, estimated_cost, status, created_at, updated_at FROM orders WHERE 1=1"
	args := []any{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if search != "" {
		query += " AND (client_name LIKE ? OR phone LIKE ?)"
		like := "%" + search + "%"
		args = append(args, like, like)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.ClientName, &o.Phone, &o.Device, &o.Description,
			&o.EstimatedCost, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *OrderRepo) GetByID(id int64) (*models.Order, error) {
	o := &models.Order{}
	err := r.db.QueryRow(
		"SELECT id, client_name, phone, device, description, estimated_cost, status, created_at, updated_at FROM orders WHERE id = ?",
		id,
	).Scan(&o.ID, &o.ClientName, &o.Phone, &o.Device, &o.Description,
		&o.EstimatedCost, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	parts, err := r.getOrderParts(id)
	if err != nil {
		return nil, err
	}
	o.Parts = parts

	return o, nil
}

func (r *OrderRepo) getOrderParts(orderID int64) ([]models.OrderPart, error) {
	rows, err := r.db.Query(`
		SELECT op.id, op.order_id, op.part_id, p.name, op.quantity, op.created_at
		FROM order_parts op
		JOIN parts p ON p.id = op.part_id
		WHERE op.order_id = ?
		ORDER BY op.created_at DESC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []models.OrderPart
	for rows.Next() {
		var op models.OrderPart
		if err := rows.Scan(&op.ID, &op.OrderID, &op.PartID, &op.PartName, &op.Quantity, &op.CreatedAt); err != nil {
			return nil, err
		}
		parts = append(parts, op)
	}
	return parts, rows.Err()
}

func (r *OrderRepo) Create(o *models.Order) (int64, error) {
	return r.CreateWithStatus(o, models.StatusAccepted)
}

func (r *OrderRepo) CreateWithStatus(o *models.Order, status models.OrderStatus) (int64, error) {
	res, err := r.db.Exec(
		"INSERT INTO orders (client_name, phone, device, description, estimated_cost, status) VALUES (?, ?, ?, ?, ?, ?)",
		o.ClientName, o.Phone, o.Device, o.Description, o.EstimatedCost, status,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *OrderRepo) Delete(id int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM order_parts WHERE order_id = ?", id); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM orders WHERE id = ?", id); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *OrderRepo) UpdateStatus(id int64, status models.OrderStatus) error {
	valid := false
	for _, s := range models.OrderStatuses {
		if s == status {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid status: %s", status)
	}

	_, err := r.db.Exec(
		"UPDATE orders SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	)
	return err
}

// StatusCounts returns count per status for the dashboard badge.
func (r *OrderRepo) StatusCounts() (map[string]int, error) {
	rows, err := r.db.Query("SELECT status, COUNT(*) FROM orders GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := map[string]int{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}

	// Also compute total
	total := 0
	for _, v := range counts {
		total += v
	}
	counts["total"] = total

	return counts, rows.Err()
}

// AllStatuses returns a display-friendly version of all statuses with labels.
func AllStatusLabels() []struct {
	Value string
	Label string
} {
	result := []struct {
		Value string
		Label string
	}{}
	for _, s := range models.OrderStatuses {
		result = append(result, struct {
			Value string
			Label string
		}{Value: string(s), Label: s.Label()})
	}
	return result
}

// SearchTokens splits a search query into normalized tokens.
func SearchTokens(q string) []string {
	return strings.Fields(strings.ToLower(strings.TrimSpace(q)))
}
