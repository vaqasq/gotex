package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"time"
	"unicode"

	"golang.org/x/term"
)

// GLOBAL VARIABLES

type Editor struct {
	Lines    [][]rune
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

	currentRowIndex       int
	currentWithinRowIndex int

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

// Represents the largest cursorX value possible for the current editor line
var furthestRight int

const (
	// ANSI TERMINAL ESCAPE CODES
	clearScreen   = "\x1b[2J"
	cursorTopLeft = "\x1b[H"

	// UNICODE DECIMAL TRANSLATON
	ESCAPE_SEQUENCE = 27
	EXIT            = 3  // Control C
	PAGE_UP         = 26 // Control Z
	PAGE_DOWN       = 24 // Control X
	PAGING_VALUE    = 20
	TAB             = 9
	BACKSPACE       = 127
	ENTER           = 13

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

	E.cursorY = 1
	E.cursorX = 1
	E.currentRow = 1

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

	E.pageUpCounter++
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

	E.pageDownCounter++
	E.pageUpCounter = 0
}

func visibleLines() [][]rune {
	if len(E.Lines) == 0 {
		return nil
	}

	start := max(0, E.rowOffset)
	end := min(len(E.Lines), start+E.config.screenRows)
	return E.Lines[start:end]
}

func displayFileContents() {
	for _, line := range visibleLines() {
		fmt.Printf("%s\r\n", string(line))
	}
	displayStatusBar()
}

func displayStatusBar() {
	fmt.Printf("\033[34m%s\033[0m \033[32m%d Lines\033[0m \033[33m%v\033[0m \033[31m%s\033[0m "+
		"currRow: %d rowOffset %d cursX %d cursY %d currWithinRowIndex %d",
		E.fileName, len(E.Lines), time.Now().Format(time.Kitchen),
		"Ctrl+C to Quit", E.currentRow, E.rowOffset, E.cursorX, E.cursorY, E.currentWithinRowIndex)
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

	updateCurrentWithinRowIndex()

	// Relies on currentRow to be updated
	// Inefficient because not every key input is an arrow movement up or down
	// Will not worry about this optimization for now
	checkBounds()

	cleanupScreen()
	displayFileContents()
	displayCursor()
}

func updateCurrentWithinRowIndex() {

	// Need to calculate the index that E.cursorX is actually hovering over because tabs mess this up.
	E.currentWithinRowIndex = E.cursorX - 1

	for _, value := range E.Lines[E.currentRowIndex] {
		if value == '\t' {
			E.currentWithinRowIndex -= 7
		}
	}

}

// This function makes sure that if a user if arrowing up and down
// that their cursor stays in the bounds of the E.Lines row
func checkBounds() {

	E.currentRowIndex = E.currentRow - 1

	// No need to add extra 1 for E.cursorX 1-indexing because len() is total runes, not 0-indexed.
	// Add 1 so that the user can go 1 further beyond the current text. Will need to be reflected in insertion code
	furthestRight = len(E.Lines[E.currentRowIndex]) + 1

	// Accounting for tabs
	for _, val := range E.Lines[E.currentRowIndex] {
		if val == '\t' {
			furthestRight += TAB_WIDTH - 1 // Already 1 space in there
		}
	}

	// checksBounds for a normal line that does not contain a tab
	if E.cursorX > furthestRight {

		E.cursorX = furthestRight

		// checkBounds for when the line begins with a tab
	} else if len(E.Lines[E.currentRowIndex]) > 0 && E.Lines[E.currentRowIndex][0] == '\t' && E.cursorX < TAB_WIDTH {

		E.cursorX = TAB_WIDTH
		E.currentWithinRowIndex = 0

	}
}

func insertRune(b rune) {

	index := E.cursorX - 1

	// E.currentRow-1 is the index for the current text line in the editor, []rune
	E.Lines[E.currentRowIndex] = slices.Insert(E.Lines[E.currentRowIndex], index, b)

	E.cursorX++
}

func tab() {
	E.Lines[E.currentRowIndex] = slices.Insert(E.Lines[E.currentRowIndex], E.currentWithinRowIndex, '\t')
	E.cursorX += TAB_WIDTH
	E.currentWithinRowIndex++
}

func backspace() {

	if E.currentWithinRowIndex < len(E.Lines[E.currentRowIndex]) {
		E.Lines[E.currentRowIndex] = slices.Delete(E.Lines[E.currentRowIndex], E.currentWithinRowIndex, E.currentWithinRowIndex+1)
	}

	if len(E.Lines[E.currentRowIndex]) == 0 {
		if E.cursorY != 1 {
			E.cursorY--
		}
		E.Lines = slices.Delete(E.Lines, E.currentRowIndex, E.currentRowIndex+1)
	}

	if E.currentWithinRowIndex == 0 && E.Lines[E.currentRowIndex][0] == '\t' {
		E.Lines[E.currentRowIndex] = slices.Delete(E.Lines[E.currentRowIndex], 0, 1)
		E.cursorX = 1
	}

}

func enter() {

	// If the enter is the at end of the current line
	if E.cursorX == len(E.Lines[E.currentRowIndex])+1 {
		E.Lines = slices.Insert(E.Lines, E.currentRowIndex+1, []rune{})

		// Otherwise in the midst of a given line
	} else {
		// use append to force a mem copy.
		remainingText := append([]rune(nil), E.Lines[E.currentRowIndex][E.currentWithinRowIndex:]...)
		E.Lines = slices.Insert(E.Lines, E.currentRowIndex+1, remainingText)
		E.Lines[E.currentRowIndex] = E.Lines[E.currentRowIndex][:E.currentWithinRowIndex]

	}
	fmt.Printf("Printing! --> %d\n", E.Lines[E.currentRowIndex+1])

	E.currentWithinRowIndex = 0
	E.cursorX = 1
	E.cursorY++

}

func moveCursor(input byte) {
	switch input {
	case UP_ARROW:
		if E.cursorY > 1 {
			E.cursorY--
		} else if E.rowOffset > 0 {
			E.rowOffset--
		}

	case DOWN_ARROW:

		// Prevents out of bounds errors by setting a limit.
		// Removes possibility of indexing out of E.Lines
		if E.rowOffset+E.cursorY >= len(E.Lines) {
			return
		}

		if E.cursorY < E.config.screenRows {
			E.cursorY++
		} else {
			E.rowOffset++
		}

	case RIGHT_ARROW:
		if E.cursorX != furthestRight {
			E.currentWithinRowIndex++
			if E.Lines[E.currentRowIndex][E.currentWithinRowIndex] == '\t' {
				E.cursorX += TAB_WIDTH
			} else {
				E.cursorX++
			}
		}
	case LEFT_ARROW:
		if E.cursorX > 1 && E.currentWithinRowIndex > 0 {
			E.currentWithinRowIndex--
			if E.Lines[E.currentRowIndex][E.currentWithinRowIndex] == '\t' {
				E.cursorX -= TAB_WIDTH
			} else {
				E.cursorX--
			}
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

	for scanner.Scan() {
		E.Lines = append(E.Lines, []rune(scanner.Text()))
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

	// Reads a single entry at a time.
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
		fmt.Println("bye\r")
		return true
	case PAGE_UP:
		pageUp()
		return false
	case PAGE_DOWN:
		pageDown()
		return false
	case TAB:
		tab()
		return false
	case BACKSPACE:
		backspace()
		return false
	case ENTER:
		enter()
		return false
	default:
		if unicode.IsPrint(rune(b)) {
			insertRune(rune(b))
		}
		return false
	}

	// for esc sequence case
	//fmt.Printf("%v pressed!\r\n", string(b))
	return false
}
