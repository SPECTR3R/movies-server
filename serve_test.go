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

func newGetYearRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/movie/%s", name), nil)
	return req
}

func newGetMoviesRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/movies", nil)
	return req
}

func newPostMovieRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/movie/", nil)
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

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func assertMovies(t testing.TB, got, want []Movie) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGETMovieYear(t *testing.T) {
	wantedMovies := []Movie{
		{1, "El aro", 1995},
		{2, "Candyman", 1992},
	}

	store := MockMovieStore{wantedMovies}

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

		store := MockMovieStore{wantedMovies}
		server := NewMovieServer(&store)

		request := newGetMoviesRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getMoviesFromResponse(t, response.Body)
		assertContentType(t, response, jsonContentType)
		assertMovies(t, got, wantedMovies)
	})
}

func getMoviesFromResponse(t testing.TB, body io.Reader) (movies []Movie) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&movies)
	if err != nil {
		t.Fatalf("Unable to parse response from server '%s' into slice of Movie, '%v'", body, err)
	}
	return movies
}

func TestPostMovie(t *testing.T) {
	store := MockMovieStore{
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
