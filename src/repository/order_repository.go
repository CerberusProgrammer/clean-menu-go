package repository

import (
	"database/sql"
	"time"

	"sazardev.clean-menu-go/src/models"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrder(order models.Order) (int, error) {
	query := `INSERT INTO orders (table_id, user_id, status, notes, payment_method, created_at, updated_at, total_amount, discount, tax) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	err := r.DB.QueryRow(query, order.TableID, order.UserID, order.Status, order.Notes, order.PaymentMethod, time.Now(), time.Now(), order.TotalAmount, order.Discount, order.Tax).Scan(&order.ID)
	if err != nil {
		return 0, err
	}

	for _, item := range order.Items {
		item.OrderID = order.ID
		err = r.CreateOrderItem(item)
		if err != nil {
			return 0, err
		}
	}

	return order.ID, nil
}

func (r *OrderRepository) CreateOrderItem(item models.OrderItem) error {
	query := `INSERT INTO order_items (order_id, menu_id, quantity, price, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, item.OrderID, item.MenuID, item.Quantity, item.Price, time.Now())
	return err
}

func (r *OrderRepository) UpdateOrder(order models.Order) error {
	query := `UPDATE orders SET table_id = $1, user_id = $2, status = $3, notes = $4, payment_method = $5, updated_at = $6, total_amount = $7, discount = $8, tax = $9 WHERE id = $10`
	_, err := r.DB.Exec(query, order.TableID, order.UserID, order.Status, order.Notes, order.PaymentMethod, time.Now(), order.TotalAmount, order.Discount, order.Tax, order.ID)
	return err
}

func (r *OrderRepository) DeleteOrder(id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	rows, err := r.DB.Query(`SELECT id, table_id, user_id, status, notes, payment_method, created_at, updated_at, total_amount, discount, tax FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.TableID, &order.UserID, &order.Status, &order.Notes, &order.PaymentMethod, &order.CreatedAt, &order.UpdatedAt, &order.TotalAmount, &order.Discount, &order.Tax)
		if err != nil {
			return nil, err
		}

		order.Items, err = r.GetOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderByID(id int) (models.Order, error) {
	var order models.Order
	query := `SELECT id, table_id, user_id, status, notes, payment_method, created_at, updated_at, total_amount, discount, tax FROM orders WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&order.ID, &order.TableID, &order.UserID, &order.Status, &order.Notes, &order.PaymentMethod, &order.CreatedAt, &order.UpdatedAt, &order.TotalAmount, &order.Discount, &order.Tax)
	if err != nil {
		return order, err
	}

	order.Items, err = r.GetOrderItems(order.ID)
	return order, err
}

func (r *OrderRepository) GetOrderItems(orderID int) ([]models.OrderItem, error) {
	rows, err := r.DB.Query(`SELECT id, order_id, menu_id, quantity, price, created_at FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.MenuID, &item.Quantity, &item.Price, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
