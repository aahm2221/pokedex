package pokeapi

import "github.com/aahm2221/pokedex/internal/pokecache"

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
	Pokemon  map[string]Pokemon
}
