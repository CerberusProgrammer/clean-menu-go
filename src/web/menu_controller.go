package web

import (
	"database/sql"
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
	"sazardev.clean-menu-go/src/repository"
)

var menuRepository *repository.MenuRepository

func InitMenuRepository(db *sql.DB) {
	menuRepository = repository.NewMenuRepository(db)
}

func ListMenus(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.GetCurrentUser()
	menus, err := menuRepository.GetAllMenus()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load menus")
		return
	}

	data := struct {
		CurrentUser models.User
		Menus       []models.Menu
	}{
		CurrentUser: currentUser,
		Menus:       menus,
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
			Name:          name,
			Price:         price,
			Description:   description,
			Recipe:        r.FormValue("recipe"),
			Availability:  r.FormValue("availability") == "on",
			EstimatedTime: atoi(r.FormValue("estimated_time")),
			Ingredients:   r.Form["ingredients"],
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

		err = menuRepository.CreateMenu(menu)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to create menu", http.StatusInternalServerError)
			return
		}
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
	currentUser := auth.GetCurrentUser()

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

		menu, err := menuRepository.GetMenuByID(id)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to load menu", http.StatusInternalServerError)
			return
		}

		menu.Name = name
		menu.Price = price
		menu.Description = description
		menu.Recipe = r.FormValue("recipe")
		menu.Availability = r.FormValue("availability") == "on"
		menu.Categories = categories
		menu.EstimatedTime = atoi(r.FormValue("estimated_time"))
		menu.Ingredients = r.Form["ingredients"]
		menu.UpdatedAt = time.Now()

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

		err = menuRepository.UpdateMenu(menu)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to update menu", http.StatusInternalServerError)
			return
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

	menu, err := menuRepository.GetMenuByID(id)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load menu")
		return
	}

	data := struct {
		CurrentUser models.User
		Menu        models.Menu
	}{
		CurrentUser: currentUser,
		Menu:        menu,
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

	err = ts.ExecuteTemplate(w, "base", data)
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

	err = menuRepository.DeleteMenu(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to delete menu", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/menus", http.StatusSeeOther)
}
