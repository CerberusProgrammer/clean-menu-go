package main

import (
	"fmt"
	"net/http"

	"sazardev.clean-menu-go/src/web"
)

func main() {
	mux := http.NewServeMux()

	// src\ui\static
	fileserver := http.FileServer(http.Dir("src/ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux.HandleFunc("/", web.Home)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", mux)
}
