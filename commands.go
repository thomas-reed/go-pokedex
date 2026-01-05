package main

import (
	"fmt"
	"os"

	"github.com/thomas-reed/go-pokedex/internal/pokeapi"
)

type cliCommand struct {
	name string
	description string
	callback func(cfg *pokeapi.Config) error
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
	}
}

func commandHelp(cfg *pokeapi.Config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for cmd, cmdData := range getCmdList() {
		fmt.Printf("%s: %s\n", cmd, cmdData.description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *pokeapi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	fmt.Println()
	os.Exit(0)
	return nil
}

func commandMap(cfg *pokeapi.Config) error {
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

func commandMapB(cfg *pokeapi.Config) error {
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