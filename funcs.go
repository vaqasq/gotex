package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

//func refreshScreen() { fmt.Print("\x1b[2J") }

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
