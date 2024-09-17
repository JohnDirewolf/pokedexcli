package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	//Just leaving this as a fun title when starting program.
	fmt.Println("Hello, Pokeworld!")
	//Create a command line reader
	cmdRead := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		//Read to delimiter.
		cmd, err := cmdRead.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again! ", err)
		}
		fmt.Printf("You said: %v", cmd)
		//Trim the delimiter, I actually probably do not have to, its just more for my own learning
		cmd = strings.TrimSuffix(cmd, "\n")
		//Change the cmd to all lower case for clean switching.

		cmd = strings.ToLower(cmd)
		if cmd == "exit" {
			break
		}
	}

}
