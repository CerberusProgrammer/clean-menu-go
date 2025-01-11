package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")
	mux := http.NewServeMux()

	// src\ui\static
	fileserver := http.FileServer(http.Dir("src/ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

}
