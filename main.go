package main

import (
	"os"

	"golang.org/x/term"
)

// Ctrl + C doesn't terminate program automatically?

func main() {

	// Raw mode
	oldState := initializeRawMode()
	defer term.Restore(int(os.Stdin.Fd()), oldState)

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
