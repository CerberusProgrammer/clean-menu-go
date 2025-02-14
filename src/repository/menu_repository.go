package repository

import (
	"database/sql"

	"github.com/lib/pq"
	"sazardev.clean-menu-go/src/models"
)

type MenuRepository struct {
	DB *sql.DB
}

func NewMenuRepository(db *sql.DB) *MenuRepository {
	return &MenuRepository{DB: db}
}

func (r *MenuRepository) CreateMenu(menu models.Menu) error {
	query := `INSERT INTO menus (name, price, recipe, categories, image, description, availability, estimated_time, ingredients, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.DB.Exec(query, menu.Name, menu.Price, menu.Recipe, pq.Array(menu.Categories), menu.Image, menu.Description, menu.Availability, menu.EstimatedTime, pq.Array(menu.Ingredients), menu.CreatedBy.ID, menu.CreatedAt, menu.UpdatedAt)
	return err
}

func (r *MenuRepository) UpdateMenu(menu models.Menu) error {
	query := `UPDATE menus SET name = $1, price = $2, recipe = $3, categories = $4, image = $5, description = $6, availability = $7, estimated_time = $8, ingredients = $9, updated_at = $10 WHERE id = $11`
	_, err := r.DB.Exec(query, menu.Name, menu.Price, menu.Recipe, pq.Array(menu.Categories), menu.Image, menu.Description, menu.Availability, menu.EstimatedTime, pq.Array(menu.Ingredients), menu.UpdatedAt, menu.ID)
	return err
}

func (r *MenuRepository) DeleteMenu(id int) error {
	query := `DELETE FROM menus WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *MenuRepository) GetAllMenus() ([]models.Menu, error) {
	rows, err := r.DB.Query(`SELECT id, name, price, recipe, categories, image, description, availability, estimated_time, ingredients, created_by, created_at, updated_at FROM menus`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []models.Menu
	for rows.Next() {
		var menu models.Menu
		var createdBy int
		err := rows.Scan(&menu.ID, &menu.Name, &menu.Price, &menu.Recipe, pq.Array(&menu.Categories), &menu.Image, &menu.Description, &menu.Availability, &menu.EstimatedTime, pq.Array(&menu.Ingredients), &createdBy, &menu.CreatedAt, &menu.UpdatedAt)
		if err != nil {
			return nil, err
		}
		menu.CreatedBy = models.User{ID: createdBy}
		menus = append(menus, menu)
	}
	return menus, nil
}

func (r *MenuRepository) GetMenuByID(id int) (models.Menu, error) {
	var menu models.Menu
	var createdBy int
	query := `SELECT id, name, price, recipe, categories, image, description, availability, estimated_time, ingredients, created_by, created_at, updated_at FROM menus WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&menu.ID, &menu.Name, &menu.Price, &menu.Recipe, pq.Array(&menu.Categories), &menu.Image, &menu.Description, &menu.Availability, &menu.EstimatedTime, pq.Array(&menu.Ingredients), &createdBy, &menu.CreatedAt, &menu.UpdatedAt)
	menu.CreatedBy = models.User{ID: createdBy}
	return menu, err
}
