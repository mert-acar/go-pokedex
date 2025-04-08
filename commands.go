package main

import (
	"fmt"
	"math/rand"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func commandHelp(c *config, args ...string) error {
	help_txt := `
Welcome to the Pokedex!
Usage:

`
	for _, cmd := range getCommands() {
		help_txt += fmt.Sprintf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println(help_txt)
	return nil
}

func commandExit(c *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func runMap(direction string, c *config) error {
	var url *string
	switch direction {
	case "next":
		url = c.nextLocationsURL
	case "prev":
		url = c.prevLocationsURL
	default:
		url = nil
	}

	area, err := c.pokeapiClient.ListArea(url)
	if err != nil {
		return err
	}
	for _, location := range area.Results {
		fmt.Println(location.Name)
	}
	c.nextLocationsURL = &area.Next
	c.prevLocationsURL = &area.Previous
	return nil
}

func commandMapb(c *config, args ...string) error {
	return runMap("prev", c)
}

func commandMapf(c *config, args ...string) error {
	return runMap("next", c)
}

func commandExplore(c *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("No location is given!")
	}
	location_name := args[0]
	fmt.Printf("Exploring %s...\n", location_name)

	location, err := c.pokeapiClient.ListLocation(location_name)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, entry := range location.PokemonEncounters {
		fmt.Println(" -", entry.Pokemon.Name)
	}
	return nil
}

func commandInspect(c *config, args ...string) error {
	pokemon_name := args[0]
	pokemon, ok := c.pokedex[pokemon_name]
	if !ok {
		return fmt.Errorf("%s is not caught yet!\n", pokemon_name)
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, stat := range pokemon.Types {
		fmt.Printf(" -%s\n", stat.Type.Name)
	}
	return nil
}

func commandCatch(c *config, args ...string) error {
	pokemon_name := args[0]
	if _, ok := c.pokedex[pokemon_name]; ok {
		return fmt.Errorf("%s is already caught and in pokedex!\n", pokemon_name)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon_name)
	pokemon, err := c.pokeapiClient.ListPokemon(pokemon_name)
	if err != nil {
		return err
	}

	required := (-0.95/112)*float64(pokemon.BaseExperience) + 1
	roll := rand.Float64()
	if roll <= required {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		c.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func commandPokedex(c *config, args ...string) error {
  if len(c.pokedex) == 0 {
    return fmt.Errorf("Your pokedex is empty!")
  }

  fmt.Println("Your pokedex:")
  for key, _ := range c.pokedex {
    fmt.Printf(" -%s\n", key)
  }
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display next 20 locations in the world",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous 20 locations in the world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a given area: explore <area-name>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a given pokemon: catch <pokemon-name>",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Inspect your pokedex",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}
