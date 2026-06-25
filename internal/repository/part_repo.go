package repository

import (
	"database/sql"
	"fmt"
	"mobile-service/internal/models"
)

type PartRepo struct {
	db *sql.DB
}

func NewPartRepo(db *sql.DB) *PartRepo {
	return &PartRepo{db: db}
}

func (r *PartRepo) List() ([]models.Part, error) {
	query := `
		SELECT p.id, p.name, p.quantity, p.purchase_price, p.sell_price, p.sku, p.category_id, c.name, p.created_at 
		FROM parts p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.name ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []models.Part
	for rows.Next() {
		var p models.Part
		var catName sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.PurchasePrice, &p.SellPrice, &p.SKU, &p.CategoryID, &catName, &p.CreatedAt); err != nil {
			return nil, err
		}
		if catName.Valid {
			p.CategoryName = catName.String
		}
		parts = append(parts, p)
	}

	// Fetch models for each part
	// For performance in a real app this might be a single query with IN clause or JOIN, but for now N+1 is fine for a small catalog, or we can fetch all part_models and group them.
	// Let's just fetch all models and map them.
	pmQuery := `
		SELECT pm.part_id, dm.id, dm.name 
		FROM part_models pm
		JOIN device_models dm ON pm.model_id = dm.id
	`
	pmRows, err := r.db.Query(pmQuery)
	if err == nil {
		defer pmRows.Close()
		modelMap := make(map[int64][]models.DeviceModel)
		for pmRows.Next() {
			var partID, modID int64
			var modName string
			if err := pmRows.Scan(&partID, &modID, &modName); err == nil {
				modelMap[partID] = append(modelMap[partID], models.DeviceModel{ID: modID, Name: modName})
			}
		}
		for i := range parts {
			parts[i].Models = modelMap[parts[i].ID]
		}
	}

	return parts, nil
}

func (r *PartRepo) GetByID(id int64) (*models.Part, error) {
	query := `
		SELECT p.id, p.name, p.quantity, p.purchase_price, p.sell_price, p.sku, p.category_id, c.name, p.created_at 
		FROM parts p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?
	`
	p := &models.Part{}
	var catName sql.NullString
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Quantity, &p.PurchasePrice, &p.SellPrice, &p.SKU, &p.CategoryID, &catName, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if catName.Valid {
		p.CategoryName = catName.String
	}

	mRows, err := r.db.Query(`
		SELECT dm.id, dm.name 
		FROM part_models pm
		JOIN device_models dm ON pm.model_id = dm.id
		WHERE pm.part_id = ?
	`, id)
	if err == nil {
		defer mRows.Close()
		for mRows.Next() {
			var dm models.DeviceModel
			if err := mRows.Scan(&dm.ID, &dm.Name); err == nil {
				p.Models = append(p.Models, dm)
			}
		}
	}

	return p, nil
}

func (r *PartRepo) Create(p *models.Part, modelIDs []int64) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO parts (name, quantity, purchase_price, sell_price, sku, category_id) VALUES (?, ?, ?, ?, ?, ?)",
		p.Name, p.Quantity, p.PurchasePrice, p.SellPrice, p.SKU, p.CategoryID,
	)
	if err != nil {
		return 0, err
	}
	partID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, mid := range modelIDs {
		if _, err := tx.Exec("INSERT INTO part_models (part_id, model_id) VALUES (?, ?)", partID, mid); err != nil {
			return 0, err
		}
	}

	return partID, tx.Commit()
}

func (r *PartRepo) Update(p *models.Part, modelIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE parts SET name = ?, quantity = ?, purchase_price = ?, sell_price = ?, sku = ?, category_id = ? WHERE id = ?",
		p.Name, p.Quantity, p.PurchasePrice, p.SellPrice, p.SKU, p.CategoryID, p.ID,
	)
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM part_models WHERE part_id = ?", p.ID); err != nil {
		return err
	}

	for _, mid := range modelIDs {
		if _, err := tx.Exec("INSERT INTO part_models (part_id, model_id) VALUES (?, ?)", p.ID, mid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PartRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM parts WHERE id = ?", id)
	return err
}

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
