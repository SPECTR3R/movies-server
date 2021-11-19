package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubMovieStore struct {
	movies []Movie
}

func (s *StubMovieStore) GetMovies() []Movie {
	return s.movies
}

func (s *StubMovieStore) GetMovieYear(name string) int {
	movies := s.GetMovies()
	for _, movie := range movies {
		if movie.Name == name {
			return movie.Year
		}
	}
	return 0
}

func (s *StubMovieStore) RecordMovie(name string) {
	// s.winCalls = append(s.winCalls, name)
}

func newGetYearRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/movie/%s", name), nil)
	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func assertMovies(t testing.TB, got, want []Movie) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newMoviesRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/movies", nil)
	return req
}

func getMoviesFromResponse(t testing.TB, body io.Reader) (movies []Movie) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&movies)
	if err != nil {
		t.Fatalf("Unable to parse response from server '%s' into slice of Movie, '%v'", body, err)
	}

	return movies
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func TestGETMovieYear(t *testing.T) {
	wantedMovies := []Movie{
		{1, "El aro", 1995},
		{2, "Candyman", 1992},
	}

	store := StubMovieStore{wantedMovies}

	server := NewMovieServer(&store)

	t.Run("returns Candyman's Year", func(t *testing.T) {
		request := newGetYearRequest("El aro")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "1995"

		assertResponseBody(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns 404 on missing movie", func(t *testing.T) {
		request := newGetYearRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestGETMovies(t *testing.T) {
	t.Run("it returns the movies table as JSON", func(t *testing.T) {
		wantedMovies := []Movie{
			{1, "Cleo", 32},
			{2, "Chris", 20},
			{3, "Tiest", 14},
		}

		store := StubMovieStore{wantedMovies}
		server := NewMovieServer(&store)

		request := newMoviesRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getMoviesFromResponse(t, response.Body)
		assertContentType(t, response, jsonContentType)
		assertMovies(t, got, wantedMovies)
	})
}

func TestStoreYear(t *testing.T) {
	store := StubMovieStore{
		nil,
	}
	server := NewMovieServer(&store)

	t.Run("it returns accepted on  POST", func(t *testing.T) {
		movie := "Pepper"
		request := newPostMovieRequest(movie)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}

func newPostMovieRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/movie/", nil)
	return req
}
