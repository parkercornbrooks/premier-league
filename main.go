package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	PORT     = ":8000"
	TEAM_URL = "https://fantasy.premierleague.com/api/bootstrap-static/"
)

type Team struct {
	Code      int    `json:"code"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type TeamData struct {
	Teams []Team `json:"teams"`
}

func fetchTeamData() (TeamData, error) {
	resp, err := http.Get(TEAM_URL)
	if err != nil {
		return TeamData{}, err
	}
	defer resp.Body.Close()

	var d TeamData

	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return TeamData{}, err
	}
	return d, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	premData, err := fetchTeamData()
	if err != nil {
		fmt.Fprintf(w, "error getting team data")
	}
	data, err := json.Marshal(premData)
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
