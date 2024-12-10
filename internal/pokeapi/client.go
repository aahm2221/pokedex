package pokeapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(cfg *Config, mapb bool) ([]string, error) {
	var url string
	if mapb {
		if cfg.Previous == "" {
			return nil, errors.New("you're on the first page")
		}
		url = cfg.Previous
	} else {
		if cfg.Next == "" {
			return nil, errors.New("you've reached the end of the list")
		}
		url = cfg.Next
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response LocationAreaResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	cfg.Previous = response.Previous
	cfg.Next = response.Next

	var names []string
	for _, item := range response.Results {
		names = append(names, item.Name)
	}
	return names, nil
}
