package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/JohnDirewolf/pokedexcli/pokecache"
)

// The Cache variable we will be using
var myCache pokecache.CacheStruct

// Location is a single location. This is pulled from a pokeapi program on github.
type Location struct {
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
	GameIndices []struct {
		GameIndex  int `json:"game_index"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"game_indices"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
}

// LocationArea is a single location area. Has pokemon encounters
type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

// Pokemon Struct
type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Forms          []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height                 int           `json:"height"`
	HeldItems              []interface{} `json:"held_items"`
	ID                     int           `json:"id"`
	IsDefault              bool          `json:"is_default"`
	LocationAreaEncounters string        `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name    string `json:"name"`
	Order   int    `json:"order"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string      `json:"back_default"`
		BackFemale       interface{} `json:"back_female"`
		BackShiny        string      `json:"back_shiny"`
		BackShinyFemale  interface{} `json:"back_shiny_female"`
		FrontDefault     string      `json:"front_default"`
		FrontFemale      interface{} `json:"front_female"`
		FrontShiny       string      `json:"front_shiny"`
		FrontShinyFemale interface{} `json:"front_shiny_female"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func InitializeCache() error {
	myCache = pokecache.NewCache(5 * time.Second)
	//Currently there should not be an error just initializing the cache, but is so can add error checking.
	return nil
}

// This tales a location id or name and returns a location struct
func GetLocation(locationNameorID string) (LocationArea, error) {
	URL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v/", locationNameorID)
	var locationVar LocationArea

	if localLocationData, found := myCache.Get(URL); found {
		err := json.Unmarshal(localLocationData, &locationVar)
		if err != nil {
			return LocationArea{}, errors.New(fmt.Sprintf("Unknown Location ID: %vn", locationNameorID))
		} else {
			return locationVar, nil
		}
	} else {
		locationData, err := http.Get(URL)
		if err != nil {
			fmt.Printf("Error accessing PokeAPI: %s\n", err)
			return LocationArea{}, err
		}
		defer locationData.Body.Close()
		decoder := json.NewDecoder(locationData.Body)
		err = decoder.Decode(&locationVar)
		//Add the existing location to cache
		if locationByte, err := json.Marshal(locationVar); err != nil {
			fmt.Printf("Error Marshalling Data on Location: %v\n", err)
			return LocationArea{}, err
		} else {
			myCache.Add(URL, locationByte)
			return locationVar, nil
		}
	}
}

// This tales a pokemon id or name and returns a pokemon struct
func GetPokemon(pokemonNameorID string) (Pokemon, error) {
	URL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v/", pokemonNameorID)
	var pokemonVar Pokemon

	if localPokemonData, found := myCache.Get(URL); found {
		err := json.Unmarshal(localPokemonData, &pokemonVar)
		if err != nil {
			return Pokemon{}, errors.New(fmt.Sprintf("Unknown Pokemon ID: %vn", pokemonNameorID))
		} else {
			return pokemonVar, nil
		}
	} else {
		pokemonData, err := http.Get(URL)
		if err != nil {
			fmt.Printf("Error accessing PokeAPI: %s\n", err)
			return Pokemon{}, err
		}
		defer pokemonData.Body.Close()
		decoder := json.NewDecoder(pokemonData.Body)
		err = decoder.Decode(&pokemonVar)
		//Add the existing pokemon to cache
		if pokemonByte, err := json.Marshal(pokemonVar); err != nil {
			fmt.Printf("Error Marshalling Data on Location: %v\n", err)
			return Pokemon{}, err
		} else {
			myCache.Add(URL, pokemonByte)
			return pokemonVar, nil
		}
	}
}
