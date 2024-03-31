package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const PORT = ":8000"

type Res struct {
	Teams    []Team   `json:"teams"`
	Fixtures Fixtures `json:"fixtures"`
}

func handler(w http.ResponseWriter, r *http.Request) {
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
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
