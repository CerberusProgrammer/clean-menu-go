package repository

import (
	"database/sql"

	"sazardev.clean-menu-go/src/models"
)

type FloorRepository struct {
	DB *sql.DB
}

func NewFloorRepository(db *sql.DB) *FloorRepository {
	return &FloorRepository{DB: db}
}

func (r *FloorRepository) CreateFloor(floor models.Floor) error {
	query := `INSERT INTO floors (name, description, is_active, "order") VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, floor.Name, floor.Description, floor.IsActive, floor.Order)
	return err
}

func (r *FloorRepository) UpdateFloor(floor models.Floor) error {
	query := `UPDATE floors SET name = $1, description = $2, is_active = $3, "order" = $4 WHERE id = $5`
	_, err := r.DB.Exec(query, floor.Name, floor.Description, floor.IsActive, floor.Order, floor.ID)
	return err
}

func (r *FloorRepository) DeleteFloor(id int) error {
	query := `DELETE FROM floors WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *FloorRepository) GetAllFloors() ([]models.Floor, error) {
	rows, err := r.DB.Query(`SELECT id, name, description, is_active, "order" FROM floors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var floors []models.Floor
	for rows.Next() {
		var floor models.Floor
		err := rows.Scan(&floor.ID, &floor.Name, &floor.Description, &floor.IsActive, &floor.Order)
		if err != nil {
			return nil, err
		}
		floors = append(floors, floor)
	}
	return floors, nil
}

func (r *FloorRepository) GetFloorByID(id int) (models.Floor, error) {
	var floor models.Floor
	query := `SELECT id, name, description, is_active, "order" FROM floors WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&floor.ID, &floor.Name, &floor.Description, &floor.IsActive, &floor.Order)
	return floor, err
}
