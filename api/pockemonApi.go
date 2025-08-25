package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func MakeRequest(url string, qs string) (LocationAreaResponse, error) {
	fullUrl := fmt.Sprintf("%s%s", url, qs)
	res, err := http.Get(fullUrl)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	if res.StatusCode > 299 {
		return LocationAreaResponse{}, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}

	var resp LocationAreaResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return LocationAreaResponse{}, err
	}

	if len(resp.Results) > 0 {
		fmt.Printf("First result: %s (%s)\n", resp.Results[0].Name, resp.Results[0].URL)
	}
	return resp, nil
}

func ExploreLocation(url string, area string) (ExploreResponse, error) {
	fullUrl := fmt.Sprintf("%s%s", url, area)
	res, err := http.Get(fullUrl)
	if err != nil {
		return ExploreResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ExploreResponse{}, err
	}

	if res.StatusCode > 299 {
		return ExploreResponse{}, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}

	var resp ExploreResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return ExploreResponse{}, err
	}

	return resp, nil
}

func GetPokemon(url, name string) (Pokemon, error) {
	fullUrl := fmt.Sprintf("%s%s", url, name)
	res, err := http.Get(fullUrl)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	if res.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}

	var p Pokemon
	if err := json.Unmarshal(body, &p); err != nil {
		return Pokemon{}, err
	}

	return p, nil
}

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous *string  `json:"previous"`
	Results  []Result `json:"results"`
}

type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
