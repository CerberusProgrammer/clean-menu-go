package models

type Menu struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	Recipe     string   `json:"recipe"`
	Categories []string `json:"categories"`
	Image      string   `json:"image"`
}

var Menus []Menu
