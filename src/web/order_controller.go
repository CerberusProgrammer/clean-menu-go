package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/models"
	"sazardev.clean-menu-go/src/repository"
)

var orderRepository *repository.OrderRepository

func InitOrderRepository(db *sql.DB) {
	orderRepository = repository.NewOrderRepository(db)
}

func ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := orderRepository.GetAllOrders()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load orders")
		return
	}

	funcMap := template.FuncMap{
		"formatDate": formatDate,
	}

	data := struct {
		CurrentUser models.User
		Orders      []models.Order
	}{
		CurrentUser: auth.GetCurrentUser(),
		Orders:      orders,
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "orders.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.New("base").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func formatDate(dateStr string) string {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return date.Format("2006/01/02 15:04")
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		order := models.Order{
			TableID:       atoi(r.FormValue("table_id")),
			UserID:        atoi(r.FormValue("user_id")),
			Status:        r.FormValue("status"),
			Notes:         r.FormValue("notes"),
			PaymentMethod: r.FormValue("payment_method"),
			CreatedAt:     time.Now().Format(time.RFC3339),
			UpdatedAt:     time.Now().Format(time.RFC3339),
		}

		for i := range r.Form["menu_id[]"] {
			orderItem := models.OrderItem{
				MenuID:    atoi(r.Form["menu_id[]"][i]),
				Quantity:  atoi(r.Form["quantity[]"][i]),
				Price:     getMenuPrice(atoi(r.Form["menu_id[]"][i])),
				CreatedAt: time.Now().Format(time.RFC3339),
			}
			order.Items = append(order.Items, orderItem)
		}

		_, err := orderRepository.CreateOrder(order)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to create order", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/orders", http.StatusSeeOther)
		return
	}

	tables, err := tableRepository.GetAllTables()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to load tables", http.StatusInternalServerError)
		return
	}

	menus, err := menuRepository.GetAllMenus()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to load menus", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentUser models.User
		Tables      []models.Table
		Menus       []models.Menu
	}{
		CurrentUser: auth.GetCurrentUser(),
		Tables:      tables,
		Menus:       menus,
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "create_order.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func EditOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))
		order, err := orderRepository.GetOrderByID(id)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to load order", http.StatusInternalServerError)
			return
		}

		order.TableID = atoi(r.FormValue("table_id"))
		order.UserID = atoi(r.FormValue("user_id"))
		order.Status = r.FormValue("status")
		order.Notes = r.FormValue("notes")
		order.PaymentMethod = r.FormValue("payment_method")
		order.UpdatedAt = time.Now().Format(time.RFC3339)

		order.Items = []models.OrderItem{}
		for i := range r.Form["menu_id[]"] {
			orderItem := models.OrderItem{
				MenuID:    atoi(r.Form["menu_id[]"][i]),
				Quantity:  atoi(r.Form["quantity[]"][i]),
				Price:     getMenuPrice(atoi(r.Form["menu_id[]"][i])),
				CreatedAt: time.Now().Format(time.RFC3339),
			}
			order.Items = append(order.Items, orderItem)
		}

		err = orderRepository.UpdateOrder(order)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to update order", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/orders", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	order, err := orderRepository.GetOrderByID(id)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load order")
		return
	}

	tables, err := tableRepository.GetAllTables()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to load tables", http.StatusInternalServerError)
		return
	}

	menus, err := menuRepository.GetAllMenus()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to load menus", http.StatusInternalServerError)
		return
	}

	data := struct {
		CurrentUser models.User
		Order       models.Order
		Tables      []models.Table
		Menus       []models.Menu
	}{
		CurrentUser: auth.GetCurrentUser(),
		Order:       order,
		Tables:      tables,
		Menus:       menus,
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "edit_order.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func ViewOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	order, err := orderRepository.GetOrderByID(id)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load order")
		return
	}

	var table models.Table
	for _, t := range models.Tables {
		if t.ID == order.TableID {
			table = t
			break
		}
	}

	user, err := userRepository.GetUserByID(order.UserID)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load user")
		return
	}

	funcMap := template.FuncMap{
		"getMenuName": getMenuName,
	}

	data := struct {
		CurrentUser models.User
		Order       models.Order
		Table       models.Table
		User        models.User
	}{
		CurrentUser: auth.GetCurrentUser(),
		Order:       order,
		Table:       table,
		User:        user,
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "view_order.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.New("base").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	err := orderRepository.DeleteOrder(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to delete order", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/orders", http.StatusSeeOther)
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func getMenuPrice(menuID int) float64 {
	for _, menu := range models.Menus {
		if menu.ID == menuID {
			return menu.Price
		}
	}
	return 0
}

func getMenuName(menuID int) string {
	for _, menu := range models.Menus {
		if menu.ID == menuID {
			return menu.Name
		}
	}
	return ""
}
