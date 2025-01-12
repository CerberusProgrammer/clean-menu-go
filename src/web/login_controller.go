package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
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

	err = ts.ExecuteTemplate(w, "base", nil)

	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}
