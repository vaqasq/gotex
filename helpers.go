package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

// GLOBAL VARIABLES

type Editor struct {
	Lines []string

	// trying out struct embedding. I isn't actually super useful here, but it is fun.
	config editorConfig

	// Cursor 2D positioning
	// Don't forget that the terminal is 1-indexed
	cursorX int
	cursorY int
}

type editorConfig struct {

	// Terminal data
	screenRows    int
	screenColumns int
}

var E Editor // initalize editorConfig. Could use struct literal instead

const (
	// ANSI TERMINAL ESCAPE CODES
	clearScreen   = "\x1b[2J"
	cursorTopLeft = "\x1b[H"

	// UNICODE TRANSLATON
	ESCAPE_SEQUENCE = 27
	EXIT            = 3

	UP_ARROW    = 'C'
	DOWN_ARROW  = 'D'
	RIGHT_ARROW = 'B'
	LEFT_ARROW  = 'A'
)

// SET IMPLEMENTATIONS

// This allows for O(1) look ups
var arrowSet = map[byte]struct{}{
	UP_ARROW:    {},
	DOWN_ARROW:  {},
	RIGHT_ARROW: {},
	LEFT_ARROW:  {},
}

// FUNCTIONS

func initializeRawMode() *term.State {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Panicf("Error in initializing raw mode: %v", err)
	}

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Panicf("Error in  size of the terminal: %v", err)
	}

	E.config = editorConfig{
		screenColumns: width,
		screenRows:    height,
	}

	return oldState
}

func displayFileContents() {
	for _, line := range E.Lines {
		fmt.Print(line)
	}
}

func displayCursor() {
	fmt.Printf("\x1b[%d;%dH", E.cursorX+1, E.cursorY+1) // The terminal is 1-indexed, so I add 1 to match up.
}

func cleanupScreen() {
	fmt.Print(clearScreen)
	fmt.Print(cursorTopLeft)
}

func refreshScreen() {
	cleanupScreen()
	displayFileContents()
	displayCursor()
}

func moveCursor(input byte) {
	switch input {
	case UP_ARROW:
		if E.cursorX != 0 {
			E.cursorY -= 1
		}
	case DOWN_ARROW:
		if E.cursorY != E.config.screenRows-2 {
			E.cursorY += 1
		}
	case RIGHT_ARROW:
		if E.cursorX != E.config.screenColumns-1 {
			E.cursorX += 1
		}
	case LEFT_ARROW:
		if E.cursorX != 0 {
			E.cursorX -= 1
		}
	}
}

func homePage() {
	fmt.Print("\nWelcome to Gotex! To edit a file, provide the name of the file as the flag when running Gotex in the command line!\r\n\n")
}

func parseArgs() *os.File {

	if len(os.Args) == 1 {
		homePage()
		return nil
	}

	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Panicf("Opening File Failed: %v", err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		E.Lines = append(E.Lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Panicf("Error in bufio.scanner: %v", err)
	}

	return file

}

func run() {

	for {

		exitEditor := processKeyPress()
		if exitEditor {
			cleanupScreen()
			return
		}
		refreshScreen()

	}

}

func processKeyPress() (exit bool) {

	// Bufio's implementation greatly reduces sys calls, requests 4kb of memory.
	reader := bufio.NewReader(os.Stdin)

	// Reads a single entry at a time
	b, err := reader.ReadByte()
	if err != nil {
		log.Panicf("Error when reading in bytes: %v", err)
	}

	switch b {
	case ESCAPE_SEQUENCE: // if the byte is an escape sequence
		nextBytes, _ := reader.Peek(2)
		// if the byte is a CSI (Control Sequence Introducer) "["
		if nextBytes[0] == '[' {
			// if the input is an arrow key
			if _, isArrowKey := arrowSet[nextBytes[1]]; isArrowKey {
				moveCursor(nextBytes[1])
				reader.Discard(2)
				return false
			}
		}
	case EXIT:
		fmt.Println("Ctrl + C was pressed!\r")
		return true
	default:
		fmt.Printf("%v pressed!\r\n", string(b))
		return false
	}

	// for esc sequence case
	fmt.Printf("%v pressed!\r\n", string(b))
	return false

}
