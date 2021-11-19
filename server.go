package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

type Movie struct {
	ID   int
	Name string
	Year int
}

type MovieStore interface {
	GetMovieYear(name string) int
	RecordMovie(name string)
	GetMovies() []Movie
}

type MovieServer struct {
	store MovieStore
	http.Handler
}

func NewMovieServer(store MovieStore) *MovieServer {
	p := new(MovieServer)
	p.store = store

	router := http.NewServeMux()
	router.Handle("/movies", http.HandlerFunc(p.moviesHandler))
	router.Handle("/movie/", http.HandlerFunc(p.movieHandler))

	p.Handler = router
	return p
}

func (p *MovieServer) getYear(w http.ResponseWriter, movie string) {
	Year := p.store.GetMovieYear(movie)
	if Year == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, p.store.GetMovieYear(movie))
}

func (p *MovieServer) postMovie(w http.ResponseWriter, movie string) {
	p.store.RecordMovie(movie)
	w.WriteHeader(http.StatusAccepted)
}

func (p *MovieServer) moviesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetMovies())
}

func (p *MovieServer) movieHandler(w http.ResponseWriter, r *http.Request) {
	movie := strings.TrimPrefix(r.URL.Path, "/movie/")
	switch r.Method {
	case http.MethodPost:
		p.postMovie(w, movie)
	case http.MethodGet:
		p.getYear(w, movie)
	}
}
