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

func ListTables(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"GetColorStatus": models.GetColorStatus,
	}

	data := struct {
		Tables []models.Table
	}{
		Tables: models.Tables,
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "tables.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
		filepath.Join("src", "ui", "components", "shape-table.component.html"),
	}

	ts, err := template.New("base").Funcs(funcMap).ParseFiles(files...)
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

func ViewTable(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var table models.Table
	for _, t := range models.Tables {
		if t.ID == id {
			table = t
			break
		}
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "view_table.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", table)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func CreateTable(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		capacity, _ := strconv.Atoi(r.FormValue("capacity"))
		table := models.Table{
			ID:       len(models.Tables) + 1,
			Number:   r.FormValue("number"),
			Name:     r.FormValue("name"),
			Capacity: capacity,
			Shape:    r.FormValue("shape"),
			IsActive: r.FormValue("is_active") == "on",
			Status:   r.FormValue("status"),
		}
		models.Tables = append(models.Tables, table)
		http.Redirect(w, r, "/tables", http.StatusSeeOther)
		return
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "create_table.tmpl.html"),
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

func EditTable(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		id, _ := strconv.Atoi(r.FormValue("id"))
		capacity, _ := strconv.Atoi(r.FormValue("capacity"))
		for i, table := range models.Tables {
			if table.ID == id {
				models.Tables[i].Number = r.FormValue("number")
				models.Tables[i].Name = r.FormValue("name")
				models.Tables[i].Capacity = capacity
				models.Tables[i].Shape = r.FormValue("shape")
				models.Tables[i].IsActive = r.FormValue("is_active") == "on"
				models.Tables[i].Status = r.FormValue("status")
				break
			}
		}
		http.Redirect(w, r, "/tables", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var table models.Table
	for _, t := range models.Tables {
		if t.ID == id {
			table = t
			break
		}
	}

	files := []string{
		filepath.Join("src", "ui", "pages", "edit_table.tmpl.html"),
		filepath.Join("src", "ui", "layouts", "layout.tmpl.html"),
		filepath.Join("src", "ui", "components", "nav.component.html"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	err = ts.ExecuteTemplate(w, "base", table)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for i, table := range models.Tables {
		if table.ID == id {
			models.Tables = append(models.Tables[:i], models.Tables[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/tables", http.StatusSeeOther)
}
