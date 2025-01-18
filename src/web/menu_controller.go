package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"sazardev.clean-menu-go/src/models"
)

func ListMenus(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join("src", "ui", "pages", "menus.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", models.Menus)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		menu := models.Menu{
			ID:       len(models.Menus) + 1,
			Name:     r.FormValue("name"),
			Price:    price,
			Recipe:   r.FormValue("recipe"),
			Category: r.FormValue("category"),
		}
		models.Menus = append(models.Menus, menu)
		http.Redirect(w, r, "/menus", http.StatusSeeOther)
		return
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "create_menu.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func EditMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		for i, menu := range models.Menus {
			if menu.ID == id {
				models.Menus[i].Name = r.FormValue("name")
				models.Menus[i].Price = price
				models.Menus[i].Recipe = r.FormValue("recipe")
				models.Menus[i].Category = r.FormValue("category")
				break
			}
		}
		http.Redirect(w, r, "/menus", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var menu models.Menu
	for _, m := range models.Menus {
		if m.ID == id {
			menu = m
			break
		}
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "edit_menu.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", menu)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for i, menu := range models.Menus {
		if menu.ID == id {
			models.Menus = append(models.Menus[:i], models.Menus[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/menus", http.StatusSeeOther)
}
