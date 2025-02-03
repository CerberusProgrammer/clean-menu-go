package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/models"
	"sazardev.clean-menu-go/src/repository"
)

var userRepository = repository.NewUserRepository(repository.DB)

func InitUserRepository(db *sql.DB) {
	userRepository = repository.NewUserRepository(db)
}

type TemplateData struct {
	CurrentUser models.User
	User        models.User
	Users       []models.User
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	currentUser := auth.GetCurrentUser()
	users, err := userRepository.GetAllUsers()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load users")
		return
	}

	data := TemplateData{
		CurrentUser: currentUser,
		Users:       users,
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
		users, err := userRepository.GetAllUsers()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to load users", http.StatusInternalServerError)
			return
		}
		for _, user := range users {
			if user.Email == email {
				http.Error(w, "Email already exists", http.StatusBadRequest)
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
			Email:    email,
			Password: password,
			Username: r.FormValue("username"),
			Name:     r.FormValue("name"),
			LastName: r.FormValue("last_name"),
			Phone:    r.FormValue("phone"),
			Role:     r.FormValue("role"),
		}

		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			uploadDir := filepath.Join("src", "ui", "static", "uploads")

			// Create the uploads directory if it doesn't exist
			err = os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
				return
			}

			filePath := filepath.Join(uploadDir, handler.Filename)
			dst, err := os.Create(filePath)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// Copy the uploaded file to the destination file
			_, err = io.Copy(dst, file)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}

			user.Image = filepath.ToSlash(filepath.Join("uploads", handler.Filename))
		}

		err = userRepository.CreateUser(user)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to create user", http.StatusInternalServerError)
			return
		}

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
		users, err := userRepository.GetAllUsers()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to load users", http.StatusInternalServerError)
			return
		}
		for _, user := range users {
			if user.Email == email && user.ID != id {
				http.Error(w, "Email already exists", http.StatusBadRequest)
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
			ID:       id,
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

			// Create the uploads directory if it doesn't exist
			err = os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
				return
			}

			filePath := filepath.Join(uploadDir, handler.Filename)
			dst, err := os.Create(filePath)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// Copy the uploaded file to the destination file
			_, err = io.Copy(dst, file)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}

			user.Image = filepath.ToSlash(filepath.Join("uploads", handler.Filename))
		}

		err = userRepository.UpdateUser(user)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to update user", http.StatusInternalServerError)
			return
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

	user, err := userRepository.GetUserByID(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to load user", http.StatusInternalServerError)
		return
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

	err = userRepository.DeleteUser(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.SetCurrentUser(models.User{})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
