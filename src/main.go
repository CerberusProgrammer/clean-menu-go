package main

import (
	"flag"
	"fmt"
	"net/http"

	"sazardev.clean-menu-go/src/auth"
	"sazardev.clean-menu-go/src/repository"
	"sazardev.clean-menu-go/src/web"
)

func getDataSourceName() string {
	userFlag := flag.String("user", "", "User name for the database connection")

	if userFlag != nil {
		fmt.Println("Error while parsing the user flag")
	}

	passwordFlag := flag.String("password", "", "Password for the database connection")

	if passwordFlag != nil {
		fmt.Println("Error while parsing the password flag")
	}

	dbNameFlag := flag.String("dbname", "", "Database name for the database connection")

	if dbNameFlag != nil {
		fmt.Println("Error while parsing the dbname flag")
	}

	port := flag.String("port", "5432", "Port for the database connection")

	if port != nil {
		fmt.Println("Error while parsing the port flag")
	}

	if userFlag == nil || passwordFlag == nil || dbNameFlag == nil {
		fmt.Println("Error while parsing the flags")
		return ""
	}

	flag.Parse()

	dataSourceName := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s port=%s", *userFlag, *dbNameFlag, *passwordFlag, *port)

	return dataSourceName
}

func main() {
	mux := http.NewServeMux()

	repository.InitDB(getDataSourceName())
	web.InitUserRepository(repository.DB)
	web.InitTableRepository(repository.DB)
	web.InitFloorRepository(repository.DB)
	web.InitMenuRepository(repository.DB)
	web.InitOrderRepository(repository.DB)

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
