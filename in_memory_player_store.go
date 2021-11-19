package main

func NewInMemoryMovieStore() *InMemoryMovieStore {
	return &InMemoryMovieStore{}
}

type InMemoryMovieStore struct{}

func (i *InMemoryMovieStore) RecordMovie(name string) {
}

func (i *InMemoryMovieStore) GetMovieYear(name string) int {
	return 0
}

func (i *InMemoryMovieStore) GetMovies() []Movie {
	var movies []Movie
	return movies
}
