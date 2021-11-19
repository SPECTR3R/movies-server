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
	RecordWin(name string)
	GetLeague() []Movie
}

type MovieServer struct {
	store MovieStore
	http.Handler
}

func NewMovieServer(store MovieStore) *MovieServer {
	p := new(MovieServer)
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/movie/", http.HandlerFunc(p.movieHandler))

	p.Handler = router
	return p
}

func (p *MovieServer) showYear(w http.ResponseWriter, movie string) {
	Year := p.store.GetMovieYear(movie)
	if Year == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, p.store.GetMovieYear(movie))
}

func (p *MovieServer) processWin(w http.ResponseWriter, movie string) {
	p.store.RecordWin(movie)
	w.WriteHeader(http.StatusAccepted)
}

func (p *MovieServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *MovieServer) movieHandler(w http.ResponseWriter, r *http.Request) {
	movie := strings.TrimPrefix(r.URL.Path, "/movie/")
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, movie)
	case http.MethodGet:
		p.showYear(w, movie)
	}
}
