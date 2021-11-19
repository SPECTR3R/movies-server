package main

func NewInMemoryMovieStore() *InMemoryMovieStore {
	return &InMemoryMovieStore{map[string]int{}}
}

type InMemoryMovieStore struct {
	store map[string]int
}

func (i *InMemoryMovieStore) RecordWin(name string) {
	i.store[name]++
}

func (i *InMemoryMovieStore) GetMovieYear(name string) int {
	return i.store[name]
}

func (i *InMemoryMovieStore) GetLeague() []Movie {
	var league []Movie
	for name, wins := range i.store {
		id := len(league) + 1
		league = append(league, Movie{id, name, wins})
	}
	return league
}
