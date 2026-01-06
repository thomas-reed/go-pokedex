package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	
	"github.com/thomas-reed/go-pokedex/internal/pokeapi"
)

func startREPL(cfg *pokeapi.Config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		cmd := cleanInput(scanner.Text())
		if len(cmd) == 0 {
			continue
		}

		cmdName := cmd[0]
		args := []string{}
		if len(cmd) > 1 {
			args = cmd[1:]
		}
		
		command, exists := getCmdList()[cmdName]
		if exists {
			err := command.callback(cfg, args...)
			if err != nil {
				fmt.Printf("Error running %s command: %v\n\n", command.name, err)
			}
		} else {
			fmt.Println("Unknown command")
			fmt.Println()
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}