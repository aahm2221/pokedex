package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/aahm2221/pokedex/internal/commands"
	"github.com/aahm2221/pokedex/internal/pokeapi"
	"github.com/aahm2221/pokedex/internal/pokecache"
)

var cfg pokeapi.Config

func init() {
	cfg = pokeapi.Config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous: "",
		Cache:    pokecache.NewCache(5 * time.Second),
	}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			continue
		}
		if err := commands.ExecuteCommand(scanner.Text(), &cfg); err != nil {
			fmt.Println(err)
			continue
		}
	}
}
