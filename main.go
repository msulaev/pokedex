package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/api"
	"strings"
)

func main() {
	config := &configUrl{
		url:      "https://pokeapi.co/api/v2/location-area/",
		next:     "",
		previous: "",
	}

	commands := map[string]*cliCommand{
		"exit": {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"help": {name: "help", description: "Displays a help message", callback: commandHelp},
		"map":  {name: "map", description: "Map a Pokemon", config: config},
		"mapb": {name: "mapb", description: "Map previous page", config: config},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		command, ok := commands[strings.ToLower(input[0])]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		switch command.name {
		case "map":
			fetchAndPrint(command, command.config.url)
		case "mapb":
			if command.config.previous == "" {
				fmt.Println("You're on the first page")
				continue
			}
			fetchAndPrint(command, command.config.previous)
		default:
			command.callback()
		}
	}
}

func fetchAndPrint(command *cliCommand, url string) {
	resp, err := api.MakeRequest(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	// Update config
	if resp.Next != "" {
		command.config.url = resp.Next
	}
	if resp.Previous != nil && *resp.Previous != "" {
		command.config.previous = *resp.Previous
	} else {
		command.config.previous = ""
	}

	for _, result := range resp.Results {
		fmt.Println(result.Name)
	}
}

func cleanInput(text string) []string {
	return strings.Fields(text)
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("  help: Displays a help message")
	fmt.Println("  exit: Exit the Pokedex")
	fmt.Println("  map: Show next page of locations")
	fmt.Println("  mapb: Show previous page of locations")
	fmt.Println()
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
	config      *configUrl
}

type configUrl struct {
	url      string
	next     string
	previous string
}
