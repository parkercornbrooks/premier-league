package main

import (
	"encoding/json"
	"net/http"
)

const TEAM_URL = "https://fantasy.premierleague.com/api/bootstrap-static/"

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
