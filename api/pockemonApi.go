package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func MakeRequest(url string) (LocationAreaResponse, error) {
	res, err := http.Get(url)
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
