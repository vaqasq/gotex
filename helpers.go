package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"golang.org/x/term"
)

type editorConfig struct {
	editorRows    int
	editorColumns int
}

var E editorConfig

const clearScreen = "\x1b[2J"
const cursorTopLeft = "\x1b[H"

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
	fmt.Print(cursorTopLeft)
}

func processKeyPress() (exit bool) {

	reader := bufio.NewReader(os.Stdin) // Bufio's implementation greatly reduces sys calls, requests 4kb of memory.

	b, err := reader.ReadByte() // reads a single entry at a time
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
