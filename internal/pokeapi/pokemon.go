package pokeapi

type Pokemon struct {
	Name   string
	Height int
	Weight int
	Stats  map[string]int
	Types  []string
}
