package repository

import (
	"database/sql"
	"mobile-service/internal/models"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) List() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, created_at FROM categories ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *CategoryRepo) GetByID(id int64) (*models.Category, error) {
	c := &models.Category{}
	err := r.db.QueryRow("SELECT id, name, created_at FROM categories WHERE id = ?", id).Scan(&c.ID, &c.Name, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *CategoryRepo) Create(c *models.Category) (int64, error) {
	res, err := r.db.Exec("INSERT INTO categories (name) VALUES (?)", c.Name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CategoryRepo) Update(c *models.Category) error {
	_, err := r.db.Exec("UPDATE categories SET name = ? WHERE id = ?", c.Name, c.ID)
	return err
}

func (r *CategoryRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = ?", id)
	return err
}
