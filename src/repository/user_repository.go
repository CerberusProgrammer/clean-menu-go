package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"sazardev.clean-menu-go/src/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user models.User) error {
	fields := []string{}
	values := []interface{}{}
	placeholders := []string{}

	if user.Username != "" {
		fields = append(fields, "username")
		values = append(values, user.Username)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if user.Password != "" {
		fields = append(fields, "password")
		values = append(values, user.Password)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if user.Name != "" {
		fields = append(fields, "name")
		values = append(values, user.Name)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if user.LastName != "" {
		fields = append(fields, "last_name")
		values = append(values, user.LastName)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if user.Email != "" {
		fields = append(fields, "email")
		values = append(values, user.Email)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if user.Phone != "" {
		fields = append(fields, "phone")
		values = append(values, user.Phone)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if user.Role != "" {
		fields = append(fields, "role")
		values = append(values, user.Role)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}
	if user.Image != "" {
		fields = append(fields, "image")
		values = append(values, user.Image)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}

	query := fmt.Sprintf("INSERT INTO users (%s) VALUES (%s)", strings.Join(fields, ", "), strings.Join(placeholders, ", "))
	fmt.Println("User data: ", user)
	fmt.Println("Query: ", query)
	fmt.Println("Values: ", values)

	_, err := r.DB.Exec(query, values...)
	if err != nil {
		fmt.Println("Error executing query: ", err)
	}
	return err
}

func (r *UserRepository) UpdateUser(user models.User) error {
	query := `UPDATE users SET username = $1, password = $2, name = $3, last_name = $4, email = $5, phone = $6, role = $7, image = $8 WHERE id = $9`
	_, err := r.DB.Exec(query, user.Username, user.Password, user.Name, user.LastName, user.Email, user.Phone, user.Role, user.Image, user.ID)
	return err
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	rows, err := r.DB.Query(`SELECT id, username, password, name, last_name, email, phone, role, image FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.LastName, &user.Email, &user.Phone, &user.Role, &user.Image); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password, name, last_name, email, phone, role, image FROM users WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Name, &user.LastName, &user.Email, &user.Phone, &user.Role, &user.Image)
	return user, err
}
