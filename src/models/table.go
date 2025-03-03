package models

type Table struct {
	ID       int    `json:"id"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Shape    string `json:"shape"`
	IsActive bool   `json:"is_active"`
	Status   string `json:"status"`
}

var Tables []Table

const (
	TableStatusAvailable = "available"
	TableStatusOccupied  = "occupied"
	TableStatusReserved  = "reserved"
	TableShapeCircle     = "circle"
	TableShapeSquare     = "square"
	TableShapeRectangle  = "rectangle"
)

func GetColorStatus(status string) string {
	switch status {
	case TableStatusAvailable:
		return "green"
	case TableStatusOccupied:
		return "red"
	case TableStatusReserved:
		return "orange"
	default:
		return "gray"
	}
}
