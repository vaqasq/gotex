package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

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

const clearScreen = "\x1b[2J"
const cursorTopLeft = "\x1b[H"

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

func processKeyPress() (exit bool) {

	// Bufio's implementation greatly reduces sys calls, requests 4kb of memory.
	reader := bufio.NewReader(os.Stdin)

	// Reads a single entry at a time
	b, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}

	if unicode.IsControl(rune(b)) {
		fmt.Println("A control button was pressed!\r")
		if b == 3 {
			fmt.Println("Ctrl + C was pressed!\r")
			return true
		}
		return false
	} else {
		fmt.Printf("%v pressed!\r\n", string(b))
		return false
	}

}
