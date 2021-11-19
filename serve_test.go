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
	Years    map[string]int
	winCalls []string
	league   []Movie
}

func (s *StubMovieStore) GetLeague() []Movie {
	return s.league
}

func (s *StubMovieStore) GetMovieYear(name string) int {
	Year := s.Years[name]
	return Year
}

func (s *StubMovieStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func newGetYearRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/movie/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/movie/%s", name), nil)
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

func assertLeague(t testing.TB, got, want []Movie) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Movie) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response from server '%s' into slice of Movie, '%v'", body, err)
	}

	return league
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func TestGETMovie(t *testing.T) {
	store := StubMovieStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}

	server := NewMovieServer(&store)

	t.Run("returns Pepper's Year", func(t *testing.T) {
		request := newGetYearRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertResponseBody(t, got, want)
		assertResponseBody(t, got, want)
	})
	t.Run("returns Floyd's Year", func(t *testing.T) {
		request := newGetYearRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

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

func TestStoreYear(t *testing.T) {
	store := StubMovieStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewMovieServer(&store)

	t.Run("it returns accepted on  POST", func(t *testing.T) {
		movie := "Pepper"
		request := newPostWinRequest(movie)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
		if store.winCalls[0] != "Pepper" {
			t.Errorf("got %s as the name of the winner want %s", store.winCalls[0], movie)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Movie{
			{1, "Cleo", 32},
			{2, "Chris", 20},
			{3, "Tiest", 14},
		}

		store := StubMovieStore{nil, nil, wantedLeague}
		server := NewMovieServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertContentType(t, response, jsonContentType)
		assertLeague(t, got, wantedLeague)
	})
}
