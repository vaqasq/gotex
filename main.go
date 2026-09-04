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

	saveFile(file.Name())
}
