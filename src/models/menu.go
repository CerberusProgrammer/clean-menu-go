package models

type Menu struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Recipe   string  `json:"recipe"`
	Category string  `json:"category"`
}

var Menus []Menu
