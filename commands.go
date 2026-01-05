package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/thomas-reed/go-pokedex/internal/pokeapi"
)

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
			name: "explore",
			description: "Prints the names of all pokemon in a given area",
			callback: commandExplore,
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
	return nil
}

func commandExplore(cfg *pokeapi.Config, args ...string) error {
	if len(args) == 0 {
		return errors.New("Location area needed.  Usage 'explore <area-name or id>")
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
	return nil
}