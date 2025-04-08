package main

import (
	"bufio"
	"fmt"
	"github.com/mert-acar/pokedex/internal/pokeapi"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

type config struct {
	pokedex          map[string]pokeapi.PokemonResponse
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
}

func start_repl(c *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		clean_input := cleanInput(scanner.Text())
		if len(clean_input) == 0 {
			continue
		}

		args := []string{}
		if len(clean_input) > 1 {
			args = clean_input[1:]
		}

		if command, ok := getCommands()[clean_input[0]]; ok {
			err := command.callback(c, args...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}
