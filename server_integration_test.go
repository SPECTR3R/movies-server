package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryMovieStore()
	server := NewMovieServer(store)
	player1 := "Pepper"

	player2 := "Juan"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player1))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player1))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player2))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player1))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "2")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Movie{
			{"Pepper", 2},
			{"Juan", 1},
		}
		assertLeague(t, got, want)
	})
}
