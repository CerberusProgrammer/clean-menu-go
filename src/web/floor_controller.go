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

var floorRepository *repository.FloorRepository

func InitFloorRepository(db *sql.DB) {
	floorRepository = repository.NewFloorRepository(db)
}

func ListFloors(w http.ResponseWriter, r *http.Request) {
	floors, err := floorRepository.GetAllFloors()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load floors")
		return
	}

	data := struct {
		CurrentUser models.User
		Floors      []models.Floor
	}{
		CurrentUser: auth.GetCurrentUser(),
		Floors:      floors,
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
		order, _ := strconv.Atoi(r.FormValue("order"))
		floor := models.Floor{
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
			IsActive:    r.FormValue("is_active") == "on",
			Order:       order,
		}
		err := floorRepository.CreateFloor(floor)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to create floor", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/floors", http.StatusSeeOther)
		return
	}

	data := struct {
		CurrentUser models.User
	}{
		CurrentUser: auth.GetCurrentUser(),
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

	err = ts.ExecuteTemplate(w, "base", data)
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
		order, _ := strconv.Atoi(r.FormValue("order"))
		floor := models.Floor{
			ID:          id,
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
			IsActive:    r.FormValue("is_active") == "on",
			Order:       order,
		}
		err := floorRepository.UpdateFloor(floor)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Unable to update floor", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/floors", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	floor, err := floorRepository.GetFloorByID(id)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to load floor")
		return
	}

	data := struct {
		CurrentUser models.User
		Floor       models.Floor
	}{
		CurrentUser: auth.GetCurrentUser(),
		Floor:       floor,
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

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Unable to render template")
		return
	}
}

func DeleteFloor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	err := floorRepository.DeleteFloor(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Unable to delete floor", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/floors", http.StatusSeeOther)
}
