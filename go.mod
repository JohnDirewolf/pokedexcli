module github.com/JohnDirewolf/pokedexcli

go 1.23.1
replace github.com/JohnDirewolf/pokedexcli/pokeapi/pokeapi v0.0.0 => ./pokeapi
replace github.com/JohnDirewolf/pokedexcli/pokecache/pokecache v0.0.0 => ./pokecache

require (
	github.com/JohnDirewolf/pokedexcli/pokeapi/pokeapi v0.0.0
	github.com/JohnDirewolf/pokedexcli/pokecache/pokecache v0.0.0
)
