package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/aahm2221/pokedex/internal/pokeapi"
)

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"catch": {
			name:        "catch [pokemon]",
			description: "Attempts to catch given pokemon and add to pokedex",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore [location]",
			description: "Displays the names of the Pokemon at the given location",
			callback:    commandExplore,
		},
		"inspect": {
			name:        "inspect [pokemon]",
			description: "Display details about a given caught pokemon",
			callback:    commandInspect,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas in Pokemon",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the last 20 location areas in Pokemon",
			callback:    commandMapb,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays the names of caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandExit(config *pokeapi.Config, param string) error {
	if param != "" {
		return fmt.Errorf("invalid Parameters")
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.Config, param string) error {
	if param != "" {
		return fmt.Errorf("invalid parameters")
	}
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	fmt.Println()
	return nil
}

func commandMap(config *pokeapi.Config, param string) error {
	if param != "" {
		return fmt.Errorf("invalid parameters")
	}
	locations, err := pokeapi.GetLocationAreas(config, false)
	if err != nil {
		return err
	}
	for _, item := range locations {
		fmt.Println(item)
	}
	return nil
}

func commandMapb(config *pokeapi.Config, param string) error {
	if param != "" {
		return fmt.Errorf("invalid parameters")
	}
	locations, err := pokeapi.GetLocationAreas(config, true)
	if err != nil {
		return err
	}
	for _, item := range locations {
		fmt.Println(item)
	}
	return nil
}
func commandExplore(config *pokeapi.Config, param string) error {
	if param == "" {
		return fmt.Errorf("invalid parameter: command needs location name")
	}
	fmt.Printf("Exploring %s...\n", param)
	pokemon, err := pokeapi.GetLocationPokemons(config, param)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, item := range pokemon {
		fmt.Printf(" - %s\n", item)
	}
	return nil
}

func commandCatch(config *pokeapi.Config, param string) error {
	if param == "" {
		return fmt.Errorf("invalid parameter: command needs Pokemon name")
	}
	_, exists := config.Pokemon[param]
	if exists {
		return fmt.Errorf("invalid pokemon: This Pokemon has already been caught")
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", param)
	pokemonCaught, err := pokeapi.CatchPokemon(config, param)
	if err != nil {
		return err
	}
	if pokemonCaught {
		fmt.Printf("%s was caught!\n", param)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", param)
	}
	return nil
}

func commandInspect(config *pokeapi.Config, param string) error {
	if param == "" {
		return fmt.Errorf("invalid parameter: command needs Pokemon name")
	}
	pokemon, exists := config.Pokemon[param]
	if !exists {
		return fmt.Errorf("you have not caught that pokemon")
	}
	fmt.Printf("Name: %s \nHeight: %d\nWeight: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)
	fmt.Println("Stats:")
	for key, val := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", key, val)
	}
	fmt.Println("Types:")
	for _, item := range pokemon.Types {
		fmt.Printf("  -%s\n", item)
	}
	return nil
}

func commandPokedex(config *pokeapi.Config, param string) error {
	if param != "" {
		return fmt.Errorf("invalid parameters")
	}
	if len(config.Pokemon) == 0 {
		return fmt.Errorf("no pokemon have been caught")
	}
	fmt.Println("Your Pokedex:")
	for key := range config.Pokemon {
		fmt.Printf("  -%s\n", key)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *pokeapi.Config, param string) error
}

func ExecuteCommand(input string, cfg *pokeapi.Config) error {
	splitInput := strings.Split(input, " ")
	command, exists := commands[splitInput[0]]
	if !exists {
		return fmt.Errorf("unknown command")
	}
	if len(splitInput) > 1 {
		return command.callback(cfg, splitInput[1])
	}
	return command.callback(cfg, "")
}

func GetCommands() map[string]cliCommand {
	return commands
}
