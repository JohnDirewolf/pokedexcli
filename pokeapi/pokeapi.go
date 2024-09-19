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
