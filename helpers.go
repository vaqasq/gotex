package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

//To Do: Get terminal window size for drawRows()

const clearScreen = "\x1b[2J"
const cursorTopLeft = "\x1b[H"

func cleanupScreen() {

	fmt.Print(clearScreen)
	fmt.Print(cursorTopLeft)

}

func drawRows() {

	for range 24 {
		fmt.Print("~\r\n")
	}

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
		log.Fatal(err)
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
