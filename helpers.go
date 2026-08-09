package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

// GLOBAL VARIABLES

type editorConfig struct {
	editorRows    int
	editorColumns int

	// Cursor 2D positioning
	// Don't forget that the terminal is 1-indexed
	cursorX int
	cursorY int
}

var E editorConfig

const (
	clearScreen   = "\x1b[2J"
	cursorTopLeft = "\x1b[H"

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
		panic(err)
	}

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	E = editorConfig{
		editorColumns: width,
		editorRows:    height,
	}

	return oldState
}

func drawRows() {

	for range E.editorRows {
		fmt.Print("~\r\n")
	}

}

func cleanupScreen() {

	fmt.Print(clearScreen)
	fmt.Print(cursorTopLeft)

}

func refreshScreen() {

	cleanupScreen()

	drawRows()
	fmt.Printf("\x1b[%d;%dH", E.cursorX+1, E.cursorY+1) // The terminal is 1-indexed, so I add 1 to match up.
}

func moveCursor(input byte) {

	switch input {
	case UP_ARROW:
		E.cursorY += 1
	case DOWN_ARROW:
		E.cursorY -= 1
	case RIGHT_ARROW:
		E.cursorX += 1
	case LEFT_ARROW:
		E.cursorX -= 1

	}
}

func processKeyPress() (exit bool) {

	// Bufio's implementation greatly reduces sys calls, requests 4kb of memory.
	reader := bufio.NewReader(os.Stdin)

	// Reads a single entry at a time
	b, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}

	switch b {
	case ESCAPE_SEQUENCE: // if the byte is an escape sequence
		nextBytes, _ := reader.Peek(2)
		// if the byte is a CSI (Control Sequence Introducer) "["
		if nextBytes[0] == '[' {
			// if the input is an arrow key
			if _, isArrow := arrowSet[nextBytes[1]]; isArrow {
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
