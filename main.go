package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// golox filepath.lox. get filepath from args. if empty run repl, if args[1] not empty run file from path. if >1 throw error
	args := os.Args
	if len(args) > 1 {
		log.Println("Usage: golox [script]")
		os.Exit(64)
	}
	if filePath := args[1]; filePath == "" {
		runPrompt()
	} else {
		runFile(filePath)
	}

}

func runPrompt() {
	//reads from command line returning tokens
	input := bufio.NewScanner(os.Stdin)
	// opens a loop. as long as input is not null, input.Scan() returns true.
	for input.Scan() {
		fmt.Print("> ")
		// imput.Text() returns most recently generated token from scanner
		line := input.Text()
		run(line)
	}
	// when input.Scan() returns false break loop. if err returned from input.Err() print the error to console. if nil the input has ended successfully
	if err := input.Err(); err != nil {
		fmt.Printf("read input error: %v", err)
	} else {
		println("End of input, exiting...")
	}
}

func runFile(path string) {
	//reads the entire file from the path as a byte array
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	//converts to string. allowing to use the byte array as text
	fileString := string(file)
	run(fileString)
}

func run(source string) {
	// init new scanner. NOT bufio.NewScanner, this is the scanner we are going to build not yet impleneted
	scanner = scan.NewScanner(source)
	tokens := scanner.scanTokens()

	for i := 0; i < len(tokens); i++ {
		fmt.Println(tokens[i])
	}
}
