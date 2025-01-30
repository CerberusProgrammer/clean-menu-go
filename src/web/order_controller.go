package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/models"
)

func ListOrders(w http.ResponseWriter, r *http.Request) {
	data := struct {
		CurrentUser models.User
		Orders      []models.Order
	}{
		CurrentUser: auth.GetCurrentUser(),
		Orders:      models.Orders,
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "orders.tmpl.html"),
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

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		order := models.Order{
			ID:            len(models.Orders) + 1,
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
				ID:        len(models.Orders) + 1,
				OrderID:   order.ID,
				MenuID:    atoi(r.Form["menu_id[]"][i]),
				Quantity:  atoi(r.Form["quantity[]"][i]),
				Price:     getMenuPrice(atoi(r.Form["menu_id[]"][i])),
				CreatedAt: time.Now().Format(time.RFC3339),
			}
			order.Items = append(order.Items, orderItem)
		}

		models.Orders = append(models.Orders, order)
		http.Redirect(w, r, "/orders", http.StatusSeeOther)
		return
	}

	data := struct {
		CurrentUser models.User
		Tables      []models.Table
		Menus       []models.Menu
	}{
		CurrentUser: auth.GetCurrentUser(),
		Tables:      models.Tables,
		Menus:       models.Menus,
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
		for i, order := range models.Orders {
			if order.ID == id {
				models.Orders[i].TableID = atoi(r.FormValue("table_id"))
				models.Orders[i].UserID = atoi(r.FormValue("user_id"))
				models.Orders[i].Status = r.FormValue("status")
				models.Orders[i].Notes = r.FormValue("notes")
				models.Orders[i].PaymentMethod = r.FormValue("payment_method")
				models.Orders[i].UpdatedAt = time.Now().Format(time.RFC3339)
				break
			}
		}
		http.Redirect(w, r, "/orders", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var order models.Order
	for _, o := range models.Orders {
		if o.ID == id {
			order = o
			break
		}
	}

	data := struct {
		CurrentUser models.User
		Order       models.Order
		Tables      []models.Table
		Menus       []models.Menu
	}{
		CurrentUser: auth.GetCurrentUser(),
		Order:       order,
		Tables:      models.Tables,
		Menus:       models.Menus,
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
	var order models.Order
	for _, o := range models.Orders {
		if o.ID == id {
			order = o
			break
		}
	}

	var table models.Table
	for _, t := range models.Tables {
		if t.ID == order.TableID {
			table = t
			break
		}
	}

	var user models.User
	for _, u := range models.Users {
		if u.ID == order.UserID {
			user = u
			break
		}
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

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for i, order := range models.Orders {
		if order.ID == id {
			models.Orders = append(models.Orders[:i], models.Orders[i+1:]...)
			break
		}
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
