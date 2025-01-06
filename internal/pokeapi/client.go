package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type LocationPokemonResponse struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func processLocationAreaResponse(response LocationAreaResponse, cfg *Config) []string {
	cfg.Previous = response.Previous
	cfg.Next = response.Next

	var names []string
	for _, item := range response.Results {
		names = append(names, item.Name)
	}
	return names
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

	if cachedData, exists := cfg.Cache.Get(url); exists {
		var response LocationAreaResponse
		err := json.Unmarshal(cachedData, &response)
		if err != nil {
			return nil, err
		}

		return processLocationAreaResponse(response, cfg), nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cfg.Cache.Add(url, body)

	var response LocationAreaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return processLocationAreaResponse(response, cfg), nil
}

func processLocationPokemonResponse(response LocationPokemonResponse, cfg *Config) []string {

	var names []string
	for _, item := range response.PokemonEncounters {
		names = append(names, item.Pokemon.Name)
	}
	return names
}

func GetLocationPokemons(cfg *Config, name string) ([]string, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + name + "/"

	if cachedData, exists := cfg.Cache.Get(url); exists {
		var response LocationPokemonResponse
		err := json.Unmarshal(cachedData, &response)
		if err != nil {
			return nil, err
		}

		return processLocationPokemonResponse(response, cfg), nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid location: %s", name)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cfg.Cache.Add(url, body)

	var response LocationPokemonResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return processLocationPokemonResponse(response, cfg), nil
}
