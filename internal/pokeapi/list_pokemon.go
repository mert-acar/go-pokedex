package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) ListPokemon(pokemon string) (PokemonResponse, error) {
	url := baseURL + "/pokemon/" + pokemon

	var locations PokemonResponse
	if val, ok := c.cache.Get(url); ok {
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return PokemonResponse{}, err
		}
		return locations, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokemonResponse{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return PokemonResponse{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return PokemonResponse{}, err
	}

	if err := json.Unmarshal(data, &locations); err != nil {
		return PokemonResponse{}, err
	}

	c.cache.Add(url, data)
	return locations, nil
}
