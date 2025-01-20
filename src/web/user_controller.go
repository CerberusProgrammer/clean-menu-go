package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/models"
)

type TemplateData struct {
	CurrentUser models.User
	User        models.User
	Users       []models.User
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.GetCurrentUser()
	data := TemplateData{
		CurrentUser: currentUser,
		Users:       models.Users,
	}

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

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.GetCurrentUser()
	if currentUser.Role != models.ADMINISTRATOR {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Verificar que el correo electrónico no se repita
		for _, user := range models.Users {
			if user.Email == email {
				http.Error(w, "Email already in use", http.StatusBadRequest)
				return
			}
		}

		// Verificar que el correo electrónico sea válido
		_, err = mail.ParseAddress(email)
		if err != nil {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}

		user := models.User{
			ID:       len(models.Users) + 1,
			Email:    email,
			Password: password,
			Username: r.FormValue("username"),
			Name:     r.FormValue("name"),
			Phone:    r.FormValue("phone"),
			Role:     r.FormValue("role"),
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
			user.Image = filepath.ToSlash(filepath.Join("uploads", handler.Filename))
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

		// Solo el administrador o el propio usuario pueden editar
		if currentUser.Role != models.ADMINISTRATOR && currentUser.ID != id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Verificar que el correo electrónico no se repita
		for _, user := range models.Users {
			if user.Email == email && user.ID != id {
				http.Error(w, "Email already in use", http.StatusBadRequest)
				return
			}
		}

		// Verificar que el correo electrónico sea válido
		_, err = mail.ParseAddress(email)
		if err != nil {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}

		for i, user := range models.Users {
			if user.ID == id {
				models.Users[i].Email = email
				models.Users[i].Password = password
				models.Users[i].Username = r.FormValue("username")
				models.Users[i].Name = r.FormValue("name")
				models.Users[i].Phone = r.FormValue("phone")

				// Solo el administrador puede cambiar el rol
				if currentUser.Role == models.ADMINISTRATOR {
					models.Users[i].Role = r.FormValue("role")
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
					models.Users[i].Image = filepath.ToSlash(filepath.Join("uploads", handler.Filename))
				}
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

	// Solo el administrador o el propio usuario pueden editar
	if currentUser.Role != models.ADMINISTRATOR && currentUser.ID != id {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var user models.User
	for _, u := range models.Users {
		if u.ID == id {
			user = u
			break
		}
	}

	data := TemplateData{
		CurrentUser: currentUser,
		User:        user,
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

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.GetCurrentUser()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// No permitir que un usuario se elimine a sí mismo
	if currentUser.ID == id {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Solo el administrador puede eliminar usuarios
	if currentUser.Role != models.ADMINISTRATOR {
		http.Error(w, "Forbidden", http.StatusForbidden)
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

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.SetCurrentUser(models.User{})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
