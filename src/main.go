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
	mux.Handle("/menus", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "waiter", "chef")(http.HandlerFunc(web.ListMenus))))
	mux.Handle("/menus/create", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "chef")(http.HandlerFunc(web.CreateMenu))))
	mux.Handle("/menus/edit", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "chef")(http.HandlerFunc(web.EditMenu))))
	mux.Handle("/menus/delete", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "chef")(http.HandlerFunc(web.DeleteMenu))))

	// User routes
	mux.Handle("/users", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.ListUsers))))
	mux.Handle("/users/create", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.CreateUser))))
	mux.Handle("/users/edit", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.EditUser))))
	mux.Handle("/users/delete", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.DeleteUser))))

	// Floor routes
	mux.Handle("/floors", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.ListFloors))))
	mux.Handle("/floors/create", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.CreateFloor))))
	mux.Handle("/floors/edit", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.EditFloor))))
	mux.Handle("/floors/delete", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.DeleteFloor))))

	// Table routes
	mux.Handle("/tables", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "waiter")(http.HandlerFunc(web.ListTables))))
	mux.Handle("/tables/view", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "waiter")(http.HandlerFunc(web.ViewTable))))
	mux.Handle("/tables/create", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.CreateTable))))
	mux.Handle("/tables/edit", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.EditTable))))
	mux.Handle("/tables/delete", auth.AuthMiddleware(auth.RoleMiddleware("administrator")(http.HandlerFunc(web.DeleteTable))))

	// Order routes
	mux.Handle("/orders", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "waiter", "chef")(http.HandlerFunc(web.ListOrders))))
	mux.Handle("/orders/create", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "waiter")(http.HandlerFunc(web.CreateOrder))))
	mux.Handle("/orders/edit", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "waiter")(http.HandlerFunc(web.EditOrder))))
	mux.Handle("/orders/delete", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "waiter")(http.HandlerFunc(web.DeleteOrder))))
	mux.Handle("/orders/view", auth.AuthMiddleware(auth.RoleMiddleware("administrator", "waiter", "chef")(http.HandlerFunc(web.ViewOrder))))

	// Logout route
	mux.HandleFunc("/logout", web.Logout)

	fmt.Println("Server is running on port 8080")
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
