package models

import "time"

type Menu struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	Recipe        string    `json:"recipe"`
	Categories    []string  `json:"categories"`
	Image         string    `json:"image"`
	Description   string    `json:"description"`
	Availability  bool      `json:"availability"`
	EstimatedTime int       `json:"estimated_time"`
	Ingredients   []string  `json:"ingredients"`
	CreatedBy     User      `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var Menus []Menu
