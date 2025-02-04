package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/models"
	"sazardev.clean-menu-go/src/repository"
)

var tableRepository *repository.TableRepository

func InitTableRepository(db *sql.DB) {
	tableRepository = repository.NewTableRepository(db)
}

func ListTables(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"GetColorStatus": models.GetColorStatus,
	}

	tables, err := tableRepository.GetAllTables()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load tables")
		return
	}

	data := struct {
		CurrentUser models.User
		Tables      []models.Table
	}{
		CurrentUser: auth.GetCurrentUser(),
		Tables:      tables,
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
	table, err := tableRepository.GetTableByID(id)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load table")
		return
	}

	data := struct {
		CurrentUser models.User
		Table       models.Table
	}{
		CurrentUser: auth.GetCurrentUser(),
		Table:       table,
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

	err = ts.ExecuteTemplate(w, "base", data)
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
			Number:   r.FormValue("number"),
			Name:     r.FormValue("name"),
			Capacity: capacity,
			Shape:    r.FormValue("shape"),
			IsActive: r.FormValue("is_active") == "on",
			Status:   r.FormValue("status"),
		}
		err := tableRepository.CreateTable(table)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to create table", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/tables", http.StatusSeeOther)
		return
	}

	data := struct {
		CurrentUser models.User
	}{
		CurrentUser: auth.GetCurrentUser(),
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

	err = ts.ExecuteTemplate(w, "base", data)
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
		table := models.Table{
			ID:       id,
			Number:   r.FormValue("number"),
			Name:     r.FormValue("name"),
			Capacity: capacity,
			Shape:    r.FormValue("shape"),
			IsActive: r.FormValue("is_active") == "on",
			Status:   r.FormValue("status"),
		}
		err := tableRepository.UpdateTable(table)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to update table", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/tables", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	table, err := tableRepository.GetTableByID(id)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load table")
		return
	}

	data := struct {
		CurrentUser models.User
		Table       models.Table
	}{
		CurrentUser: auth.GetCurrentUser(),
		Table:       table,
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

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	err := tableRepository.DeleteTable(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to delete table", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/tables", http.StatusSeeOther)
}
