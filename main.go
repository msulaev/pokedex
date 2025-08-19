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
	comands := map[string]*cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Map a Pokemon",
			config:      config,
		},
		"mapb": {
			name:        "mapb",
			description: "Map a Pokemon",
			config:      config,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		words := cleanInput(line)
		if len(words) == 0 {
			continue
		}
		comand, ok := comands[strings.ToLower(words[0])]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if comand.name == "map" {
			resp, err := api.MakeRequest(comand.config.url)
			if err != nil {
				fmt.Println("Error making request:", err)
			}
			comand.config.url = resp.Next
			if resp.Previous != nil && *resp.Previous != "" {
				comand.config.previous = *resp.Previous
			} else {
				comand.config.previous = ""
			}
			for i := 0; i < len(resp.Results); i++ {
				fmt.Println(resp.Results[i].Name)
			}
		} else if comand.name == "mapb" {
			if comand.config.previous == "" {
				fmt.Println("you're on the first page")
				continue
			}
			resp, _ := api.MakeRequest(comand.config.previous)
			for i := 0; i < len(resp.Results); i++ {
				fmt.Println(resp.Results[i].Name)
			}
		} else {
			comand.callback()
		}
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
