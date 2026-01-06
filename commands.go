package main

import (
	"errors"
	"fmt"
	"os"
	"math/rand"

	"github.com/thomas-reed/go-pokedex/internal/pokeapi"
)

const userLevel = 25

type cliCommand struct {
	name string
	description string
	callback func(cfg *pokeapi.Config, args ...string) error
}

func getCmdList() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name: "help",
			description: "List all commands for the Pokedex",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "Print the next page of locations",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Print the previous page of locations",
			callback: commandMapB,
		},
		"explore": {
			name: "explore <area_name | id>",
			description: "Prints the names of all pokemon in a given area",
			callback: commandExplore,
		},
		"catch": {
			name: "catch <pokemon_name>",
			description: "Attempts to catch a given pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect <pokemon_name>",
			description: "Prints stats for a given pokemon if it's been caught",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "Lists the names of all pokemon that have been caught",
			callback: commandPokedex,
		},
	}
}

func commandHelp(cfg *pokeapi.Config, args ...string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmdData := range getCmdList() {
		fmt.Printf("%s: %s\n", cmdData.name, cmdData.description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *pokeapi.Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	fmt.Println()
	os.Exit(0)
	return nil
}

func commandMap(cfg *pokeapi.Config, args ...string) error {
	mapRes, err := cfg.PokeapiClient.ListLocations(cfg.NextLocationURL)
	if err != nil {
		return err
	}
	cfg.NextLocationURL = mapRes.Next
	cfg.PreviousLocationURL = mapRes.Previous

	for _, location := range mapRes.Results {
		fmt.Println(location.Name)
	}
	fmt.Println()
	return nil
}

func commandMapB(cfg *pokeapi.Config, args ...string) error {
	mapRes, err := cfg.PokeapiClient.ListLocations(cfg.PreviousLocationURL)
	if err != nil {
		return err
	}
	cfg.NextLocationURL = mapRes.Next
	cfg.PreviousLocationURL = mapRes.Previous

	for _, location := range mapRes.Results {
		fmt.Println(location.Name)
	}
	fmt.Println()
	return nil
}

func commandExplore(cfg *pokeapi.Config, args ...string) error {
	if len(args) == 0 {
		return errors.New("Location area needed.  Usage: explore <area-name or id>")
	}
	exploreRes, err := cfg.PokeapiClient.ExploreLocation(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", exploreRes.Name)
	fmt.Println("Found Pokemon:")
	for _, encounter := range exploreRes.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	fmt.Println()
	return nil
}

func commandCatch(cfg *pokeapi.Config, args ...string) error {
	if len(args) == 0 {
		return errors.New("Pokemon name needed.  Usage: catch <pokemon-name>")
	}
	pokemonRes, err := cfg.PokeapiClient.Pokemon(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonRes.Name)
	if userLevel < rand.Intn(pokemonRes.BaseExperience) {
		fmt.Printf("%s escaped!\n", pokemonRes.Name)
		fmt.Println()
		return nil
	}
	cfg.Pokedex[pokemonRes.Name] = pokemonRes
	fmt.Printf("%s was caught!\n", pokemonRes.Name)
	fmt.Println("You may now inspect it with the inspect command.")
	fmt.Println()
	return nil
}

func commandInspect(cfg *pokeapi.Config, args ...string) error {
	if len(args) == 0 {
		return errors.New("Pokemon name needed.  Usage: inspect <pokemon-name>")
	}
	pokemon, caught := cfg.Pokedex[args[0]]
	if !caught {
		fmt.Printf("You have not caught %s\n", args[0])
		fmt.Println()
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf(" - %s\n", pokeType.Type.Name)
	}
	fmt.Println()
	return nil
}

func commandPokedex(cfg *pokeapi.Config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.Pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	fmt.Println()
	return nil
}