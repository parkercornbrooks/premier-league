package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const FIXTURES_URL = "https://fantasy.premierleague.com/api/fixtures/"

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
	TeamH      int       `json:"team_h"`
	TeamAScore int       `json:"team_a_score"`
	TeamHScore int       `json:"team_h_score"`
	Finished   bool      `json:"finished"`
	Kickoff    time.Time `json:"kickoff_time"`
}

type Fixtures []Fixture

func fetchFixtureData() (Fixtures, error) {
	resp, err := http.Get(FIXTURES_URL)
	if err != nil {
		return Fixtures{}, fmt.Errorf("error fetching fixture data %w", err)
	}
	defer resp.Body.Close()

	var d Fixtures

	err = json.NewDecoder(resp.Body).Decode(&d)
	if err != nil {
		return Fixtures{}, fmt.Errorf("error decoding fixture data %w", err)
	}
	return d, nil
}
