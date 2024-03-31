package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	PORT         = ":8000"
	TEAM_URL     = "https://fantasy.premierleague.com/api/bootstrap-static/"
	FIXTURES_URL = "https://fantasy.premierleague.com/api/fixtures/"
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

type Kickoff time.Time

func (k *Kickoff) UnmarshalJSON(b []byte) error {
	s := string(b)
	t, err := time.Parse("2006-01-02T15:04:05Z", s)
	if err != nil {
		return err
	}
	*k = Kickoff(t)
	return nil
}

func (k Kickoff) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(k))
}

type Fixture struct {
	Id         int       `json:"id"`
	TeamA      int       `json:"team_a"`
	TeamB      int       `json:"team_b"`
	TeamAScore int       `json:"team_a_score"`
	TeamBScore int       `json:"team_b_score"`
	Finished   bool      `json:"finished"`
	Kickoff    time.Time `json:"kickoff_time"`
}

type Fixtures []Fixture

func fetchFixtureData() (Fixtures, error) {
	resp, err := http.Get(FIXTURES_URL)
	if err != nil {
		return Fixtures{}, err
	}
	defer resp.Body.Close()

	var d Fixtures

	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return Fixtures{}, err
	}
	return d, nil
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
	response := struct {
		Teams    []Team   `json:"teams"`
		Fixtures Fixtures `json:"fixtures"`
	}{
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
