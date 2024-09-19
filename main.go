package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JohnDirewolf/pokedexcli/pokeapi"
	"github.com/JohnDirewolf/pokedexcli/pokecache"
)

// Commands is basically a constant, it is in this structure as this what is directed by the project specs.
// I have this as a global variable as it is accessed by various functions as if it was a constant variable.
var commands map[string]cliCommand

// Location Start ID for Map and MapB, defined here so both functions can access it.
var locationStartID int = 1

// The Cache variable we will be using
var myCache pokecache.CacheStruct

func getKnownCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
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
			description: "Return 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back 20 locations",
			callback:    commandMapB,
		},
	}
	return commands
}

func commandHelp() error {
	fmt.Println("-----Pokedex Help-----")
	fmt.Println("Commands:")
	for command := range commands {
		fmt.Printf("     %s: %s\n", commands[command].name, commands[command].description)
	}
	fmt.Println("     ")

	return nil
}

func commandExit() error {
	//This is just for the cliCommand structure and if we need to do something before closing. Now it just returns and the main func breaks.
	return nil
}

func commandMap() error {

	for locationID := locationStartID; locationID < locationStartID+20; locationID++ {
		URL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", locationID)
		var locationVar pokeapi.Location

		if localLocationData, found := myCache.Get(URL); found {
			err := json.Unmarshal(localLocationData, &locationVar)
			if err != nil {
				fmt.Printf("Unknown Location ID: %d\n", locationID)
			} else {
				fmt.Println(locationVar.Name)
			}
		} else {
			locationData, err := http.Get(URL)
			defer locationData.Body.Close()
			if err != nil {
				fmt.Printf("Error accessing PokeAPI: %s\n", err)
				return err
			}
			decoder := json.NewDecoder(locationData.Body)
			err = decoder.Decode(&locationVar)
			//Add the existing location to cache
			if locationByte, err := json.Marshal(locationVar); err != nil {
				fmt.Printf("Error Marshalling Data on Location: %v\n", err)
			} else {
				myCache.Add(URL, locationByte)
			}

			if err != nil {
				fmt.Printf("Unknown Location ID: %d\n", locationID)
			} else {
				fmt.Println(locationVar.Name)
			}
		}
	}
	locationStartID = locationStartID + 20

	return nil
}

func commandMapB() error {

	//fmt.Println("Start command MapB")
	if locationStartID == 1 {
		//This should only trigger if the user tries to use MapB before Map on start of program.
		fmt.Println("Already at start of Index, please use Map.")
	} else {
		//fmt.Printf("Going back: Start location: %v", locationStartID)
		//So Map B goes back to the previous set of locations, if we have only printed the first set of locations we just reset to the start.
		locationStartID = locationStartID - 40
		if locationStartID < 1 {
			locationStartID = 1
		}
	}
	commandMap()
	return nil
}

func main() {
	//Initialize our list of known commands
	commands = getKnownCommands()
	//Create the Cache
	myCache = pokecache.NewCache(5 * time.Second)

	fmt.Println("Welcome to the Pokedex!")
	//Create a Scanner
	cmdScan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		//Scan command line for input
		cmdScan.Scan()
		//Get command and I put it to all lower to help in processing the commands.
		cmd := strings.ToLower(cmdScan.Text())
		//Check if it is a valid command, if so process it based on the fuction in the structure or alert user the command is invalid
		//fmt.Printf("Your command: %v\n", cmd)
		if _, exists := commands[cmd]; exists {
			commands[cmd].callback()
		} else {
			fmt.Println("Unknown command: Type 'help' for valid commands.")
			fmt.Println("")
		}
		//if the cmmand is exit then we just end program.
		if cmd == "exit" {
			break
		}
	}

}
