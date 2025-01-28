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

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		CurrentUser models.User
	}{
		CurrentUser: auth.GetCurrentUser(),
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "home.tmpl.html"),
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
