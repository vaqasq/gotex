package main

import (
	"fmt"
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

	// Trying to figure out how \t looks
	for _, str := range E.Lines {
		fmt.Println(str[0], "\r")
	}
	run()

}
