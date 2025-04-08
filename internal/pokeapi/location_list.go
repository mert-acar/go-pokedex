package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) ListLocation(location string) (LocationResponse, error) {
	url := baseURL + "/location-area/" + location

	var pokemon LocationResponse
	if val, ok := c.cache.Get(url); ok {
		err := json.Unmarshal(val, &pokemon)
		if err != nil {
			return LocationResponse{}, err
		}
		return pokemon, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationResponse{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return LocationResponse{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return LocationResponse{}, err
	}

	if err := json.Unmarshal(data, &pokemon); err != nil {
		return LocationResponse{}, err
	}

	c.cache.Add(url, data)
	return pokemon, nil
}
