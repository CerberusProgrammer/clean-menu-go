package repository

import (
	"time"

	"sazardev.clean-menu-go/src/models"
)

func GetOrders() ([]models.Order, error) {
	rows, err := models.DB.Query("SELECT id, table_id, user_id, status, notes, payment_method, created_at, updated_at, total_amount, discount, tax FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.TableID, &order.UserID, &order.Status, &order.Notes, &order.PaymentMethod, &order.CreatedAt, &order.UpdatedAt, &order.TotalAmount, &order.Discount, &order.Tax); err != nil {
			return nil, err
		}
		order.Items, _ = GetOrderItems(order.ID)
		orders = append(orders, order)
	}
	return orders, nil
}

func GetOrderById(id int) (models.Order, error) {
	var order models.Order
	err := models.DB.QueryRow("SELECT id, table_id, user_id, status, notes, payment_method, created_at, updated_at, total_amount, discount, tax FROM orders WHERE id = $1", id).Scan(&order.ID, &order.TableID, &order.UserID, &order.Status, &order.Notes, &order.PaymentMethod, &order.CreatedAt, &order.UpdatedAt, &order.TotalAmount, &order.Discount, &order.Tax)
	if err != nil {
		return order, err
	}
	order.Items, _ = GetOrderItems(order.ID)
	return order, nil
}

func CreateOrder(order *models.Order) error {
	err := models.DB.QueryRow("INSERT INTO orders (table_id, user_id, status, notes, payment_method, created_at, updated_at, total_amount, discount, tax) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", order.TableID, order.UserID, order.Status, order.Notes, order.PaymentMethod, time.Now(), time.Now(), order.TotalAmount, order.Discount, order.Tax).Scan(&order.ID)
	if err != nil {
		return err
	}
	for _, item := range order.Items {
		item.OrderID = order.ID
		if err := CreateOrderItem(&item); err != nil {
			return err
		}
	}
	return nil
}

func UpdateOrder(order *models.Order) error {
	_, err := models.DB.Exec("UPDATE orders SET table_id = $1, user_id = $2, status = $3, notes = $4, payment_method = $5, updated_at = $6, total_amount = $7, discount = $8, tax = $9 WHERE id = $10", order.TableID, order.UserID, order.Status, order.Notes, order.PaymentMethod, time.Now(), order.TotalAmount, order.Discount, order.Tax, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOrder(id int) error {
	_, err := models.DB.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func GetOrderItems(orderID int) ([]models.OrderItem, error) {
	rows, err := models.DB.Query("SELECT id, order_id, menu_id, quantity, price, created_at FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.MenuID, &item.Quantity, &item.Price, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func CreateOrderItem(item *models.OrderItem) error {
	err := models.DB.QueryRow("INSERT INTO order_items (order_id, menu_id, quantity, price, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id", item.OrderID, item.MenuID, item.Quantity, item.Price, time.Now()).Scan(&item.ID)
	if err != nil {
		return err
	}
	return nil
}
