package commands

import (
	"fmt"
	"os"

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
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
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
	}
}

func commandExit(config *pokeapi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.Config) error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandMap(config *pokeapi.Config) error {
	locations, err := pokeapi.GetLocationAreas(config, false)
	if err != nil {
		return err
	}
	for _, item := range locations {
		fmt.Println(item)
	}
	return nil
}

func commandMapb(config *pokeapi.Config) error {
	locations, err := pokeapi.GetLocationAreas(config, true)
	if err != nil {
		return err
	}
	for _, item := range locations {
		fmt.Println(item)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *pokeapi.Config) error
}

func ExecuteCommand(input string, cfg *pokeapi.Config) error {
	command, exists := commands[input]
	if !exists {
		return fmt.Errorf("Unknown command")
	}
	return command.callback(cfg)
}

func GetCommands() map[string]cliCommand {
	return commands
}
