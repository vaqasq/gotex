package main

import (
	"os"

	"golang.org/x/term"
)

// Ctrl + C doesn't terminate program automatically?

func main() {

	// Raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Main loop
	for {

		exit := processKeyPress()
		if exit {
			cleanupScreen()
			return
		}

	}

}
