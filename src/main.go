package main

import (
	"fmt"
	"net/http"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/web"
)

func main() {
	mux := http.NewServeMux()

	// src\ui\static
	fileserver := http.FileServer(http.Dir("src/ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/", web.Login)
	mux.Handle("/home", auth.AuthMiddleware(http.HandlerFunc(web.Home)))
	mux.Handle("/register", http.HandlerFunc(web.Register))

	// Menu routes
	mux.Handle("/menus", auth.AuthMiddleware(http.HandlerFunc(web.ListMenus)))
	mux.Handle("/menus/create", auth.AuthMiddleware(http.HandlerFunc(web.CreateMenu)))
	mux.Handle("/menus/edit", auth.AuthMiddleware(http.HandlerFunc(web.EditMenu)))
	mux.Handle("/menus/delete", auth.AuthMiddleware(http.HandlerFunc(web.DeleteMenu)))

	// User routes
	mux.Handle("/users", auth.AuthMiddleware(http.HandlerFunc(web.ListUsers)))
	mux.Handle("/users/create", auth.AuthMiddleware(http.HandlerFunc(web.CreateUser)))
	mux.Handle("/users/edit", auth.AuthMiddleware(http.HandlerFunc(web.EditUser)))
	mux.Handle("/users/delete", auth.AuthMiddleware(http.HandlerFunc(web.DeleteUser)))

	// Floor routes
	mux.Handle("/floors", auth.AuthMiddleware(http.HandlerFunc(web.ListFloors)))
	mux.Handle("/floors/create", auth.AuthMiddleware(http.HandlerFunc(web.CreateFloor)))
	mux.Handle("/floors/edit", auth.AuthMiddleware(http.HandlerFunc(web.EditFloor)))
	mux.Handle("/floors/delete", auth.AuthMiddleware(http.HandlerFunc(web.DeleteFloor)))

	// Table routes
	mux.Handle("/tables", auth.AuthMiddleware(http.HandlerFunc(web.ListTables)))
	mux.Handle("/tables/create", auth.AuthMiddleware(http.HandlerFunc(web.CreateTable)))
	mux.Handle("/tables/edit", auth.AuthMiddleware(http.HandlerFunc(web.EditTable)))
	mux.Handle("/tables/delete", auth.AuthMiddleware(http.HandlerFunc(web.DeleteTable)))

	// Logout route
	mux.HandleFunc("/logout", web.Logout)

	fmt.Println("Server is running on port 8080")
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
