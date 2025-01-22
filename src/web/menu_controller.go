package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/models"
)

func ListMenus(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.GetCurrentUser()
	data := struct {
		CurrentUser models.User
		Menus       []models.Menu
	}{
		CurrentUser: currentUser,
		Menus:       models.Menus,
	}

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

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.GetCurrentUser()
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}
		description := r.FormValue("description")
		categories := r.Form["categories[]"]

		menu := models.Menu{
			ID:            len(models.Menus) + 1,
			Name:          name,
			Price:         price,
			Description:   description,
			Recipe:        "",
			Availability:  false,
			EstimatedTime: 0,
			Ingredients:   []string{},
			Categories:    categories,
			CreatedBy:     currentUser,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			uploadDir := filepath.Join("src", "ui", "static", "uploads")
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				os.Mkdir(uploadDir, os.ModePerm)
			}
			filePath := filepath.Join(uploadDir, handler.Filename)
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
			menu.Image = filepath.ToSlash(filepath.Join("uploads", handler.Filename))
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
		err := r.ParseMultipartForm(10 << 20)
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

		name := r.FormValue("name")
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}
		description := r.FormValue("description")
		categories := r.Form["categories"]

		for i, menu := range models.Menus {
			if menu.ID == id {
				models.Menus[i].Name = name
				models.Menus[i].Price = price
				models.Menus[i].Description = description
				models.Menus[i].Recipe = r.FormValue("recipe")
				models.Menus[i].Availability = r.FormValue("availability") == "on"
				models.Menus[i].Categories = categories
				models.Menus[i].UpdatedAt = time.Now()

				file, handler, err := r.FormFile("image")
				if err == nil {
					defer file.Close()
					uploadDir := filepath.Join("src", "ui", "static", "uploads")
					if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
						os.Mkdir(uploadDir, os.ModePerm)
					}
					filePath := filepath.Join(uploadDir, handler.Filename)
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
					models.Menus[i].Image = filepath.ToSlash(filepath.Join("uploads", handler.Filename))
				}
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
