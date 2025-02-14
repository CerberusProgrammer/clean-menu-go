package repository

import (
	"database/sql"

	"sazardev.clean-menu-go/src/models"
)

type TableRepository struct {
	DB *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{DB: db}
}

func (r *TableRepository) CreateTable(table models.Table) error {
	query := `INSERT INTO tables (number, name, capacity, shape, is_active, status) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.DB.Exec(query, table.Number, table.Name, table.Capacity, table.Shape, table.IsActive, table.Status)
	return err
}

func (r *TableRepository) UpdateTable(table models.Table) error {
	query := `UPDATE tables SET number = $1, name = $2, capacity = $3, shape = $4, is_active = $5, status = $6 WHERE id = $7`
	_, err := r.DB.Exec(query, table.Number, table.Name, table.Capacity, table.Shape, table.IsActive, table.Status, table.ID)
	return err
}

func (r *TableRepository) DeleteTable(id int) error {
	query := `DELETE FROM tables WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *TableRepository) GetAllTables() ([]models.Table, error) {
	rows, err := r.DB.Query(`SELECT id, number, name, capacity, shape, is_active, status FROM tables`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []models.Table
	for rows.Next() {
		var table models.Table
		err := rows.Scan(&table.ID, &table.Number, &table.Name, &table.Capacity, &table.Shape, &table.IsActive, &table.Status)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (r *TableRepository) GetTableByID(id int) (models.Table, error) {
	var table models.Table
	query := `SELECT id, number, name, capacity, shape, is_active, status FROM tables WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&table.ID, &table.Number, &table.Name, &table.Capacity, &table.Shape, &table.IsActive, &table.Status)
	return table, err
}
