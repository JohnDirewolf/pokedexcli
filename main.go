package main

import (
	"bufio"
	//"encoding/json"
	"fmt"
	//"net/http"
	"errors"
	"os"
	"strconv"
	"strings"

	//"time"

	"github.com/JohnDirewolf/pokedexcli/pokeapi"
	//"github.com/JohnDirewolf/pokedexcli/pokecache"
)

// Commands is basically a constant, it is in this structure as this what is directed by the project specs.
// I have this as a global variable as it is accessed by various functions as if it was a constant variable.
var commands map[string]cliCommand

// Location Start ID for Map and MapB, defined here so both functions can access it.
var locationStartID int = 1

func commandHelp(notused string) error {
	fmt.Println("-----Pokedex Help-----")
	fmt.Println("Commands:")
	for command := range commands {
		fmt.Printf("     %s: %s\n", commands[command].name, commands[command].description)
	}
	fmt.Println("     ")

	return nil
}

func commandExit(notused string) error {
	//This is just for the cliCommand structure and if we need to do something before closing. Now it just returns and the main func breaks.
	return nil
}

func commandMap(notused string) error {

	for locationID := locationStartID; locationID < locationStartID+20; locationID++ {
		if locationVar, err := pokeapi.GetLocation(strconv.Itoa(locationID)); err == nil {
			fmt.Println(locationVar.Name)
		} else {
			fmt.Println(err)
		}
	}
	locationStartID = locationStartID + 20

	return nil
}

func commandMapB(notused string) error {

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
	commandMap("")
	return nil
}

func commandExplore(areaName string) error {
	var locationVar pokeapi.LocationArea
	var err error

	//first check if we have an area.
	if areaName == "" {
		fmt.Println("No Area Provided. Proper usage: <Explore area-name>")
		return errors.New("No area provided")
	} else {
		if locationVar, err = pokeapi.GetLocation(areaName); err == nil {
			//So err is nil if the api returned a value or it returned NO value on not found so we need to check that.
			if locationVar.Name == "" {
				fmt.Println("Area not found! Check your location name!")
				return errors.New("Area not found")
			}
		} else {
			fmt.Printf("Could not retrieve that location! Error Reported: %v\n", err)
			return err
		}
	}
	//Got through all the problems so we have a location variable!
	fmt.Printf("Exploring %s...\n", locationVar.Name)
	for i := range locationVar.PokemonEncounters {
		fmt.Printf(" - %s\n", locationVar.PokemonEncounters[i].Pokemon.Name)
	}

	return nil
}

func main() {
	//Initialize our list of known commands
	commands = getKnownCommands()
	//Intialize the Cache
	pokeapi.InitializeCache()

	fmt.Println("Welcome to the Pokedex!")
	//Create a Scanner
	cmdScan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		//Scan command line for input
		cmdScan.Scan()
		//Get command and I put it to all lower to help in processing the commands.
		cmdLine := strings.Split(strings.ToLower(cmdScan.Text()), " ")
		//Split cmdLine into the actually command "cmd" and parameter "par" if command takes a parameter eg explore
		cmd := cmdLine[0]
		//by defalut par is empty string
		var par string = ""
		if len(cmdLine) > 1 {
			par = cmdLine[1]
		}
		//Check if it is a valid command, if so process it based on the fuction in the structure or alert user the command is invalid
		//fmt.Printf("Your command: %v\n", cmd)
		if _, exists := commands[cmd]; exists {
			commands[cmd].callback(par)
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
