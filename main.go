package main

import (
	"os"

	"golang.org/x/term"
)

func main() {

	oldState := initializeRawMode()
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	file := parseArgs()
	if file == nil {
		return
	}
	defer file.Close()

	run()

	// Current issue: Since ex.txt is quite large, it breaks the editor.
	// You need to be able keep track of which line you are at, and make sure you are not going past E.config.screenRows
	// This will allow you to go downwards. This is the next chapter in antirez's book.

}
