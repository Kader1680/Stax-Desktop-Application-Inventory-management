package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"stax-desktop/internal/domain"
)

type CategoryRepo struct {
	DB *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{DB: db}
}

func (r *CategoryRepo) Create(ctx context.Context, c *domain.Category) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `INSERT INTO categories (name) VALUES (?)`, c.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return 0, fmt.Errorf("une catégorie avec ce nom existe déjà")
		}
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, created_at FROM categories ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Category
	for rows.Next() {
		var c domain.Category
		var createdAt string
		if err := rows.Scan(&c.ID, &c.Name, &createdAt); err != nil {
			return nil, err
		}
		c.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM categories WHERE id = ?`, id)
	return err
}
