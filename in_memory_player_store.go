package main

func NewInMemoryMovieStore() *InMemoryMovieStore {
	return &InMemoryMovieStore{map[string]int{}}
}

type InMemoryMovieStore struct {
	store map[string]int
}

func (i *InMemoryMovieStore) RecordMovie(name string) {
	i.store[name]++
}

func (i *InMemoryMovieStore) GetMovieYear(name string) int {
	return i.store[name]
}

func (i *InMemoryMovieStore) GetMovies() []Movie {
	var movies []Movie
	for name, wins := range i.store {
		id := len(movies) + 1
		movies = append(movies, Movie{id, name, wins})
	}
	return movies
}
