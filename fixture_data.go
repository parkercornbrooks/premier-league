package main

import (
	"encoding/json"
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
