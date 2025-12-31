package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

func startREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		cmd := cleanInput(scanner.Text())
		if len(cmd) == 0 {
			continue
		}

		cmdName := cmd[0]
		command, exists := getCmdList()[cmdName]
		if exists {
			err := command.callback()
			if err != nil {
				fmt.Printf("Error running %s command: %v", command.name, err)
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