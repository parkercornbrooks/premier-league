package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const PORT = ":8000"

type Res struct {
	Teams    []Team   `json:"teams"`
	Fixtures Fixtures `json:"fixtures"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	premData, err := fetchTeamData()
	if err != nil {
		fmt.Fprintf(w, "error getting team data")
	}
	fixtures, err := fetchFixtureData()
	if err != nil {
		fmt.Fprintf(w, "error getting fixture stats")
	}
	response := Res{
		Teams:    premData.Teams,
		Fixtures: fixtures,
	}
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(w, "error marshalling data")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	duration := time.Since(start)
	fmt.Println("regular handler took", duration)
}

func concHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

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
		fmt.Fprintf(w, "error fetching data")
	}
	response := Res{
		Teams:    premData.Teams,
		Fixtures: fixtures,
	}
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(w, "error marshalling data")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	duration := time.Since(start)
	fmt.Println("concurrent handler took", duration)
}

func main() {
	http.HandleFunc("/sync", handler)
	http.HandleFunc("/conc", concHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	log.Fatal(http.ListenAndServe(PORT, nil))
}
