package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// GLOBAL VARIABLES

type Editor struct {
	Lines    []string
	fileName string

	// trying out struct embedding. It isn't actually super useful here, but it is fun.
	config editorConfig

	// Cursor 2D positioning
	// Don't forget that the terminal is 1-indexed
	cursorX int
	cursorY int

	rowOffset int

	// This is 1-indexed. Subtract to get the index in E.Lines
	currentRow int

	// These act like static variables in C!
	//
	// Example
	// If the user tries to pageDown when they have already reached the bottom of the file,
	// I will move the cursor to the bottom
	pageDownCounter int
	pageUpCounter   int
}

type editorConfig struct {

	// Terminal data
	screenRows    int
	screenColumns int
}

// initalize Editor. Could use struct literal instead
var E Editor

const (
	// ANSI TERMINAL ESCAPE CODES
	clearScreen   = "\x1b[2J"
	cursorTopLeft = "\x1b[H"

	// UNICODE TRANSLATON
	ESCAPE_SEQUENCE = 27
	EXIT            = 3  // Control C
	PAGE_UP         = 26 // Control Z
	PAGE_DOWN       = 24 // Control X
	PAGING_VALUE    = 20

	// These are preceeded by the ESCAPE_SEQUENCE + [
	UP_ARROW    = 'A'
	DOWN_ARROW  = 'B'
	RIGHT_ARROW = 'C'
	LEFT_ARROW  = 'D'

	// OTHER
	TAB_WIDTH = 8
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
		// Makes space for status bar
		screenRows: height - 1,
	}

	E.rowOffset = 0
	E.cursorY = 1
	E.cursorX = 1

	return oldState
}

// I choose 15 Lines

func pageUp() {

	// Paging Logic
	if E.rowOffset > PAGING_VALUE {
		E.rowOffset -= PAGING_VALUE
	} else {
		E.rowOffset = 0
		// Bring cursor to top of file
		if E.pageUpCounter > 1 {
			E.cursorY = 1
			E.cursorX = 1
		}
	}

	E.pageUpCounter += 1
	E.pageDownCounter = 0
}

func pageDown() {

	//Paging Logic
	linesLeft := len(E.Lines) - (E.rowOffset + E.config.screenRows)
	if linesLeft > PAGING_VALUE {
		E.rowOffset += PAGING_VALUE
	} else {
		E.rowOffset += linesLeft

		// Bring cursor to bottom of file
		if E.pageDownCounter > 1 {
			E.cursorY = E.config.screenRows
			E.cursorX = 1
		}
	}

	E.pageDownCounter += 1
	E.pageUpCounter = 0
}

func visibleLines() []string {
	if len(E.Lines) == 0 {
		return nil
	}

	start := max(0, E.rowOffset)
	end := min(len(E.Lines), start+E.config.screenRows)
	return E.Lines[start:end]
}

func displayFileContents() {
	for _, line := range visibleLines() {
		fmt.Printf("%s\r\n", line)
	}
	displayStatusBar()
}

func displayStatusBar() {
	fmt.Printf("\033[34m%s\033[0m \033[32m%d Lines\033[0m \033[33m%v\033[0m \033[31m%s\033[0m "+
		"currRow: %d rowOffset %d cursX %d cursY %d",
		E.fileName, len(E.Lines), time.Now().Format(time.Kitchen),
		"Ctrl+C to Quit", E.currentRow, E.rowOffset, E.cursorX, E.cursorY)
}

func displayCursor() {
	// The terminal is 1-indexed, so I add 1 to match up.
	fmt.Printf("\x1b[%d;%dH", E.cursorY, E.cursorX)
}

func updateCurrentRowVar() {
	E.currentRow = E.rowOffset + E.cursorY
}

func cleanupScreen() {
	fmt.Print(clearScreen)
	fmt.Print(cursorTopLeft)
}

func refreshScreen() {
	// updates MUST happen before redrawing, duh.
	updateCurrentRowVar()

	// Relies on currentRow to be updated
	// Inefficient because not every key input is an arrow movement up or down
	// Will not worry about this optimization for now
	checkBounds()

	cleanupScreen()
	displayFileContents()
	displayCursor()
}

// This function makes sure that if a user if arrowing up and down
// that their cursor stays in the bounds of the E.Lines row
func checkBounds() {
	if E.cursorX > len(E.Lines[E.currentRow-1]) {
		E.cursorX = len(E.Lines[E.currentRow-1])

		// Currently not detecting the \t properly??
		if strings.Contains(E.Lines[E.currentRow-1], "\t") {
			E.cursorX += TAB_WIDTH
		}

	}
}

func moveCursor(input byte) {
	switch input {
	case UP_ARROW:
		if E.cursorY > 1 {
			E.cursorY -= 1
		} else if E.rowOffset > 0 {
			E.rowOffset -= 1
		}

	case DOWN_ARROW:

		// Prevents out of bounds errors by setting a limit.
		// Removes possibility of indexing out of E.Lines
		if E.rowOffset+E.cursorY >= len(E.Lines) {
			return
		}

		if E.cursorY < E.config.screenRows {
			E.cursorY += 1
		} else {
			E.rowOffset += 1
		}

	case RIGHT_ARROW:
		if E.cursorX != len(E.Lines[E.currentRow-1]) {
			E.cursorX += 1
		}
	case LEFT_ARROW:
		if E.cursorX > 1 {
			E.cursorX -= 1
		}
	}
}

func homePage() {
	fmt.Print("\nWelcome to Gotex! To edit a file, provide the name of the " +
		"file as the flag when running Gotex in the command line!\r\n\n")
}

func parseArgs() *os.File {

	if len(os.Args) == 1 {
		homePage()
		return nil
	}

	E.fileName = os.Args[1]
	file, err := os.Open(E.fileName)
	if err != nil {
		log.Panicf("Opening File Failed: %v", err)
	}

	scanner := bufio.NewScanner(file)

	// First line not being printed to terminal, added this to fix it.
	E.Lines = append(E.Lines, "\n")

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

		refreshScreen()
		exitEditor := processKeyPress()
		if exitEditor {
			cleanupScreen()
			return
		}

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
	case PAGE_UP:
		pageUp()
		return false
	case PAGE_DOWN:
		pageDown()
		return false
	default:
		fmt.Printf("%v pressed!\r\n", string(b))
		return false
	}

	// for esc sequence case
	fmt.Printf("%v pressed!\r\n", string(b))
	return false
}
