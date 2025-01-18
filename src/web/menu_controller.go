package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sazardev.clean-menu-go/src/models"
)

func joinStrings(sep string, elements []string) string {
	return strings.Join(elements, sep)
}

func ListMenus(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join("src", "ui", "pages", "menus.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	funcMap := template.FuncMap{
		"join": joinStrings,
	}

	ts, err := template.New("menus.tmpl.html").Funcs(funcMap).ParseFiles(files...)
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
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}

		categories := strings.Split(r.FormValue("categories"), ",")
		menu := models.Menu{
			ID:         len(models.Menus) + 1,
			Name:       r.FormValue("name"),
			Price:      price,
			Recipe:     r.FormValue("recipe"),
			Categories: categories,
		}

		// Handle file upload
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			// Create the uploads directory if it doesn't exist
			if _, err := os.Stat("uploads"); os.IsNotExist(err) {
				os.Mkdir("uploads", os.ModePerm)
			}
			// Save the file
			filePath := filepath.Join("uploads", handler.Filename)
			dst, err := os.Create(filePath)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			if _, err := dst.ReadFrom(file); err != nil {
				log.Println(err.Error())
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			menu.Image = filePath
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

	funcMap := template.FuncMap{
		"join": joinStrings,
	}

	ts, err := template.New("create_menu.tmpl.html").Funcs(funcMap).ParseFiles(files...)
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
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}

		categories := strings.Split(r.FormValue("categories"), ",")
		for i, menu := range models.Menus {
			if menu.ID == id {
				models.Menus[i].Name = r.FormValue("name")
				models.Menus[i].Price = price
				models.Menus[i].Recipe = r.FormValue("recipe")
				models.Menus[i].Categories = categories
				break
			}
		}
		http.Redirect(w, r, "/menus", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

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

	funcMap := template.FuncMap{
		"join": joinStrings,
	}

	ts, err := template.New("edit_menu.tmpl.html").Funcs(funcMap).ParseFiles(files...)
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, menu := range models.Menus {
		if menu.ID == id {
			models.Menus = append(models.Menus[:i], models.Menus[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/menus", http.StatusSeeOther)
}
