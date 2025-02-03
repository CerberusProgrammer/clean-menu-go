package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var errorMessage, email, password string

	if r.Method == http.MethodPost {
		r.ParseForm()
		email = r.FormValue("email")
		password = r.FormValue("password")

		users, err := userRepository.GetAllUsers()
		if err != nil {
			log.Println(err.Error())
			errorMessage = "Unable to load users"
		} else {
			for _, user := range users {
				if user.Email == email && user.Password == password {
					auth.SetCurrentUser(user)
					http.Redirect(w, r, "/home", http.StatusSeeOther)
					return
				}
			}
			errorMessage = "Invalid email or password"
		}
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "login.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "focus.tmpl.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	data := struct {
		ErrorMessage string
		Email        string
		Password     string
	}{
		ErrorMessage: errorMessage,
		Email:        email,
		Password:     password,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		user := models.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Role:     models.ADMINISTRATOR,
		}

		err := userRepository.CreateUser(user)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to register user", http.StatusInternalServerError)
			return
		}

		auth.SetCurrentUser(user)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "register.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "focus.tmpl.html"),
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
