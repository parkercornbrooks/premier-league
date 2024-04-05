package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

const (
	API_URL    = "https://fantasy.premierleague.com/api/"
	FIXTURE_EP = "fixtures"
	TEAM_EP    = "bootstrap-static"
)

type PremRepository struct {
	client *http.Client
}

func (pr PremRepository) getData() (Res, error) {
	var (
		premData            TeamData
		fixtures            Fixtures
		teamErr, fixtureErr error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		teamErr = pr.fetchEndpoint(TEAM_EP, &premData)
	}()
	go func() {
		defer wg.Done()
		fixtureErr = pr.fetchEndpoint(FIXTURE_EP, &fixtures)
	}()
	wg.Wait()
	if teamErr != nil || fixtureErr != nil {
		return Res{}, fmt.Errorf("error fetching data")
	}
	data := Res{
		Teams:    premData.Teams,
		Fixtures: fixtures,
	}
	return data, nil
}

// fetch endpoint calls a given URL and attempts to write the data to a passed pointer
func (pr PremRepository) fetchEndpoint(endpoint string, d interface{}) error {
	fullURL, err := url.JoinPath(API_URL, endpoint)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return err
	}
	res, err := pr.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(d)
	if err != nil {
		return err
	}
	return nil
}
