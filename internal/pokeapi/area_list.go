package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) ListArea(pageURL *string) (LocationAreaResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	var locations LocationAreaResponse
	if val, ok := c.cache.Get(url); ok {
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return LocationAreaResponse{}, err
		}
		return locations, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	if err := json.Unmarshal(data, &locations); err != nil {
		return LocationAreaResponse{}, err
	}

	c.cache.Add(url, data)
	return locations, nil
}
