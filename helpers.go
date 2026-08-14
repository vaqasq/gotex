package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/term"
)

// GLOBAL VARIABLES

type Editor struct {
	Lines    []string
	fileName string

	// trying out struct embedding. I isn't actually super useful here, but it is fun.
	config editorConfig

	// Cursor 2D positioning
	// Don't forget that the terminal is 1-indexed
	cursorX int
	cursorY int

	rowOffset int

	currentRow int
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
		screenRows:    height - 1, // Makes space for status bar
	}

	E.rowOffset = 0
	E.cursorY = 1
	E.cursorX = 1

	return oldState
}

// I choose 15 Lines

func pageUp() {
	if E.rowOffset > PAGING_VALUE {
		E.rowOffset -= PAGING_VALUE
	} else {
		E.rowOffset = 0
	}
}

func pageDown() {
	linesLeft := len(E.Lines) - (E.rowOffset + E.config.screenRows)
	if linesLeft > PAGING_VALUE {
		E.rowOffset += PAGING_VALUE
	} else {
		E.rowOffset += linesLeft
	}
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
	fmt.Printf("\033[34m%s\033[0m \033[32m%d Lines\033[0m \033[33m%v\033[0m \033[31m%s\033[0m currRow: %d rowOffset %d cursY %d", E.fileName, len(E.Lines), time.Now().Format(time.Kitchen), "Ctrl+C to Quit", E.currentRow, E.rowOffset, E.cursorY)
}

func displayCursor() {
	fmt.Printf("\x1b[%d;%dH", E.cursorY, E.cursorX) // The terminal is 1-indexed, so I add 1 to match up.
}

func updateCurrentRow() {
	E.currentRow = E.rowOffset + E.cursorY + 1
}

func cleanupScreen() {
	fmt.Print(clearScreen)
	fmt.Print(cursorTopLeft)
}

func refreshScreen() {
	cleanupScreen()
	displayFileContents()
	displayCursor()
	updateCurrentRow()
}

func moveCursor(input byte) {
	switch input {
	case UP_ARROW:
		if E.cursorY > 0 {
			E.cursorY -= 1
		} else if E.rowOffset > 0 {
			E.rowOffset -= 1
		}
	case DOWN_ARROW:

		// Prevents out of bounds errors by setting a limit. You can only go so far beyond the file. Removes possibility of indexing out of E.Lines
		if E.rowOffset+E.cursorY >= len(E.Lines)-1 {
			return
		}

		if E.cursorY < E.config.screenRows-1 {
			E.cursorY += 1
		} else {
			E.rowOffset += 1
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

	E.fileName = os.Args[1]
	file, err := os.Open(E.fileName)
	if err != nil {
		log.Panicf("Opening File Failed: %v", err)
	}

	scanner := bufio.NewScanner(file)

	E.Lines = append(E.Lines, "\n") // First line not being printed to terminal, added this to fix it.

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
