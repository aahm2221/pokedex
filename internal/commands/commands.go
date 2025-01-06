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
		return fmt.Errorf("invalid Parameters")
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
		return fmt.Errorf("invalid Parameters")
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
		return fmt.Errorf("invalid Parameters")
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
	locations, err := pokeapi.GetLocationPokemons(config, param)
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
