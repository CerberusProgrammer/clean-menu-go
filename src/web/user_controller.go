package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

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
