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

	fmt.Println("Server is running on port 8080")
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
