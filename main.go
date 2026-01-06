package main

import (
	"time"

	"github.com/thomas-reed/go-pokedex/internal/pokeapi"
)

func main() {
	client := pokeapi.NewClient(5 * time.Second, 0)
	cfg := &pokeapi.Config{
		PokeapiClient: client,
		Pokedex: make(map[string]pokeapi.Pokemon),
	}
	startREPL(cfg)
}