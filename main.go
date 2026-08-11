package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

// Ctrl + C doesn't terminate program automatically?

func main() {

	// Raw mode
	oldState := initializeRawMode()
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	if len(os.Args) == 1 {
		homePage()
		return
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Panicf("Opening File Failed: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	fmt.Println(scanner.Text(), "\r")

	// Main loop
	for {

		exitEditor := processKeyPress()
		if exitEditor {
			cleanupScreen()
			return
		}

		refreshScreen()

	}

}
