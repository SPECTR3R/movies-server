package main

type MockMovieStore struct {
	movies []Movie
}

func (m *MockMovieStore) GetMovies() []Movie {
	return m.movies
}

func (m *MockMovieStore) GetMovieYear(name string) int {
	movies := m.GetMovies()
	for _, movie := range movies {
		if movie.Name == name {
			return movie.Year
		}
	}
	return 0
}

func (m *MockMovieStore) RecordMovie(name string) {
	// s.winCalls = append(s.winCalls, name)
}
