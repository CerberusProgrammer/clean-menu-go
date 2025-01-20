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

func ListUsers(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join("src", "ui", "pages", "users.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", models.Users)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		user := models.User{
			ID:       len(models.Users) + 1,
			Username: r.FormValue("username"),
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Role:     r.FormValue("role"),
		}

		models.Users = append(models.Users, user)
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "create_user.tmpl.html"),
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

func EditUser(w http.ResponseWriter, r *http.Request) {
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

		for i, user := range models.Users {
			if user.ID == id {
				models.Users[i].Username = r.FormValue("username")
				models.Users[i].Name = r.FormValue("name")
				models.Users[i].Email = r.FormValue("email")
				models.Users[i].Phone = r.FormValue("phone")
				models.Users[i].Role = r.FormValue("role")
				break
			}
		}
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var user models.User
	for _, u := range models.Users {
		if u.ID == id {
			user = u
			break
		}
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "edit_user.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", user)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, user := range models.Users {
		if user.ID == id {
			models.Users = append(models.Users[:i], models.Users[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
