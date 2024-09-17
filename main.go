package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Commands is basically a constant, it is in this structure as this what is directed by the project specs.
// I have this as a global variable as it is accessed by various functions as if it was a constant variable.
var commands map[string]cliCommand

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

func main() {
	//Initialize our list of known commands
	commands = getKnownCommands()

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
