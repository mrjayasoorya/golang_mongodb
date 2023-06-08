package main

import (
	"log"
	books "main/api/controller/books"

	"net/http"
)

func main() {
	//Init Router
	router := http.NewServeMux()

	router.HandleFunc("/api/books", books.GetBooks)
	router.HandleFunc("/api/requestBook", books.RequestBook)
	router.HandleFunc("/api/createLibraryCard", books.CreateLibraryCard)
	// r.HandleFunc("/api/books", createBook).Methods("POST")
	// r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	// r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// config := helper.GetConfiguration()
	dir := http.Dir("./static")

	fs := http.FileServer(dir)

	router.Handle("/", fs)

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080", router))

}
