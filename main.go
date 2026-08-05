package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func main() {

	/*
		if len(os.Args) == 1 {
			log.Fatal("No file name supplied")
		} else if len(os.Args) > 2 {
			log.Fatal("Too mang arguments supplied")
		}

		fileName := os.Args[1]

		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		fmt.Println(file.Name())
	*/

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buffer := make([]byte, 1)
	var input string
	var ch byte

	for {

		_, err := os.Stdin.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		ch = buffer[0] // utf8

		if ch == 3 {
			fmt.Println("Ctrl+C Pressed")
			return
		} else if ch == 13 {
			fmt.Println("Newline Character Pressed")
			fmt.Println(input)
			return
		} else {
			fmt.Printf("%v pressed!\n", ch)
			return
		}

		//input += string(buffer[0])

	}

}
