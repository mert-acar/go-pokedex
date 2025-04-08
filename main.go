package main

import (
	"github.com/mert-acar/pokedex/internal/pokeapi"
	"time"
)

func main() {
	cfg := &config{
		pokedex:       make(map[string]pokeapi.PokemonResponse),
		pokeapiClient: pokeapi.NewClient(5*time.Second, 5*time.Second),
	}
	start_repl(cfg)
}
