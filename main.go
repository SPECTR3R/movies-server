package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewMovieServer(&InMemoryMovieStore{})
	log.Fatal(http.ListenAndServe(":8080", server))
}
