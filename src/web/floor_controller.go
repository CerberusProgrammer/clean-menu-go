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

func ListFloors(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Floors []models.Floor
	}{
		Floors: models.Floors,
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "floors.tmpl.html"),
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

func CreateFloor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		floor := models.Floor{
			ID:          len(models.Floors) + 1,
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
			IsActive:    r.FormValue("is_active") == "on",
			Order:       len(models.Floors) + 1,
		}
		models.Floors = append(models.Floors, floor)
		http.Redirect(w, r, "/floors", http.StatusSeeOther)
		return
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "create_floor.tmpl.html"),
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

func EditFloor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))
		for i, floor := range models.Floors {
			if floor.ID == id {
				models.Floors[i].Name = r.FormValue("name")
				models.Floors[i].Description = r.FormValue("description")
				models.Floors[i].IsActive = r.FormValue("is_active") == "on"
				break
			}
		}
		http.Redirect(w, r, "/floors", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var floor models.Floor
	for _, f := range models.Floors {
		if f.ID == id {
			floor = f
			break
		}
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "edit_floor.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", floor)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func DeleteFloor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for i, floor := range models.Floors {
		if floor.ID == id {
			models.Floors = append(models.Floors[:i], models.Floors[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/floors", http.StatusSeeOther)
}
