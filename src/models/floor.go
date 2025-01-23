package models

type Floor struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
	Order       int     `json:"order"`
	Tables      []Table `json:"tables"`
}

var Floors []Floor
