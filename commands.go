package main

type cliCommand struct {
	name        string
	description string
	callback    func(arg string) error
}

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
		"explore": {
			name:        "explore",
			description: "explore [area-name or id]. See a list of all the Pokemon in the area.",
			callback:    commandExplore,
		},
	}
	return commands
}
