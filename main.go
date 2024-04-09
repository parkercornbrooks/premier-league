package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const PORT = ":8000"

type Res struct {
	Teams    []Team   `json:"teams"`
	Fixtures Fixtures `json:"fixtures"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var (
		premData            TeamData
		fixtures            Fixtures
		teamErr, fixtureErr error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		premData, teamErr = fetchTeamData()
	}()
	go func() {
		defer wg.Done()
		fixtures, fixtureErr = fetchFixtureData()
	}()
	wg.Wait()
	if teamErr != nil || fixtureErr != nil {
		http.Error(w, "could not fetch data", http.StatusNotFound)
		return
	}
	response := Res{
		Teams:    premData.Teams,
		Fixtures: fixtures,
	}
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
