package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"pokedexcli/api"
	"pokedexcli/internal/pokecache"
	"strings"
	"time"
)

var pokedex = make(map[string]api.Pokemon)

func main() {
	c := pokecache.NewCache(10 * time.Second)
	defer c.Stop()
	config := &configUrl{
		url:      "https://pokeapi.co/api/v2/location-area/",
		next:     "",
		previous: "",
	}

	commands := map[string]*cliCommand{
		"exit":    {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"help":    {name: "help", description: "Displays a help message", callback: commandHelp},
		"map":     {name: "map", description: "Map a Pokemon", config: config},
		"mapb":    {name: "mapb", description: "Map previous page", config: config},
		"explore": {name: "explore", description: "Explore a Pokemon", config: config},
		"catch":   {name: "catch", description: "Catch a Pokemon", config: config},
		"inspect": {name: "inspect", description: "Inspect a caught Pokemon", config: config},
		"pokedex": {name: "pokedex", description: "Show all caught Pokemon", config: config},
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
			fetchAndPrint(command, command.config.url, "", c)
		case "mapb":
			if command.config.previous == "" {
				fmt.Println("You're on the first page")
				continue
			}
			fetchAndPrint(command, command.config.previous, "", c)
		case "explore":
			if len(input) < 2 {
				fmt.Println("Please provide a location area name. Example: explore pastoria-city-area")
				continue
			}
			fetchAndExplore(command, command.config.url, input[1], c)
		case "catch":
			if len(input) < 2 {
				fmt.Println("Please provide a pokemon name. Example: catch pikachu")
				continue
			}
			catchPokemon(command, input[1], c)
		case "inspect":
			if len(input) < 2 {
				fmt.Println("Please provide a pokemon name. Example: inspect pidgey")
				continue
			}
			inspectPokemon(input[1])
		case "pokedex":
			showPokedex()
		default:
			command.callback()
		}
	}
}

func fetchAndPrint(command *cliCommand, url string, qs string, cache *pokecache.Cache) {
	if data, ok := cache.Get(url); ok {
		var resp api.LocationAreaResponse
		if err := json.Unmarshal(data, &resp); err == nil {
			printResults(command, &resp)
			return
		}
	}

	resp, err := api.MakeRequest(url, qs)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	raw, _ := json.Marshal(resp)
	cache.Add(url, raw)

	printResults(command, &resp)
}

func fetchAndExplore(command *cliCommand, url string, area string, cache *pokecache.Cache) {
	cacheKey := url + area
	if data, ok := cache.Get(cacheKey); ok {
		var resp api.ExploreResponse
		if err := json.Unmarshal(data, &resp); err == nil {
			printExplore(area, &resp)
			return
		}
	}

	resp, err := api.ExploreLocation(url, area)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	raw, _ := json.Marshal(resp)
	cache.Add(cacheKey, raw)

	printExplore(area, &resp)
}

func printExplore(area string, resp *api.ExploreResponse) {
	fmt.Printf("Exploring %s...\n", area)
	fmt.Println("Found Pokemon:")
	for _, p := range resp.PokemonEncounters {
		fmt.Printf(" - %s\n", p.Pokemon.Name)
	}
}

func printResults(command *cliCommand, resp *api.LocationAreaResponse) {
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

	for _, result := range resp.Results {
		fmt.Println(result.Name)
	}
}

func catchPokemon(command *cliCommand, name string, cache *pokecache.Cache) {
	cacheKey := "pokemon:" + name
	var p api.Pokemon

	if data, ok := cache.Get(cacheKey); ok {
		if err := json.Unmarshal(data, &p); err != nil {
			fmt.Println("Error decoding cached pokemon:", err)
			return
		}
	} else {
		resp, err := api.GetPokemon("https://pokeapi.co/api/v2/pokemon/", name)
		if err != nil {
			fmt.Println("Error fetching pokemon:", err)
			return
		}
		p = resp
		raw, _ := json.Marshal(p)
		cache.Add(cacheKey, raw)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	if _, ok := pokedex[p.Name]; ok {
		fmt.Printf("%s is already in your Pokedex!\n", p.Name)
		return
	}

	chance := rand.Intn(p.BaseExperience + 100) // число от 0 до base+100
	if chance > p.BaseExperience/2 {
		fmt.Printf("%s was caught!\n", p.Name)
		pokedex[p.Name] = p
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", p.Name)
	}
}

func inspectPokemon(name string) {
	p, ok := pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)

	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
}

func showPokedex() {
	if len(pokedex) == 0 {
		fmt.Println("Your Pokedex is empty. Go catch some Pokemon!")
		return
	}

	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf(" - %s\n", name)
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
