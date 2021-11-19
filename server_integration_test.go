package main

import (
	"testing"
)

func TestRecordingYearAndRetrievingThem(t *testing.T) {
	// store := NewInMemoryMovieStore()
	// server := NewMovieServer(store)
	// movie1 := "Pepper"

	// movie2 := "Juan"

	// server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(movie1))
	// server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(movie1))
	// server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(movie2))

	// t.Run("get Year", func(t *testing.T) {
	// 	response := httptest.NewRecorder()
	// 	server.ServeHTTP(response, newGetYearRequest(movie1))
	// 	assertStatus(t, response.Code, http.StatusOK)

	// 	assertResponseBody(t, response.Body.String(), "2")
	// })

	// t.Run("get movies", func(t *testing.T) {
	// 	response := httptest.NewRecorder()
	// 	server.ServeHTTP(response, newMoviesRequest())
	// 	assertStatus(t, response.Code, http.StatusOK)

	// 	got := getMoviesFromResponse(t, response.Body)
	// 	want := []Movie{
	// 		{"Pepper", 2},
	// 		{"Juan", 1},
	// 	}
	// 	assertMovies(t, got, want)
	// })
}
