package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/constwhite/golox-interpreter/interpreter"
	"github.com/constwhite/golox-interpreter/parser"
	"github.com/constwhite/golox-interpreter/resolver"
	"github.com/constwhite/golox-interpreter/scanner"
)

var hadError bool
var hadRuntimeError bool

func main() {
	// golox filepath.lox. get filepath from args. if empty run repl, if args[1] not empty run file from path. if >1 throw error
	args := os.Args
	if len(args) > 2 {
		log.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runPrompt()

	} else if len(args) == 2 {
		runFile(args[1])
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
		hadError = false
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
	if hadError {
		os.Exit(65)
	}

	if hadRuntimeError {
		os.Exit(70)
	}

}

func run(source string) {
	// init new scanner. NOT bufio.NewScanner, this is the scanner we are going to build not yet impleneted

	scanner := scanner.NewScanner(source, os.Stderr)
	tokens := scanner.ScanTokens()

	parser := parser.NewParser(tokens, os.Stderr)
	statements, HadError := parser.Parse()

	if HadError {
		return
	}
	interpreter := interpreter.NewInterpreter(os.Stderr, os.Stdout)
	resolver := resolver.NewResolver(interpreter)
	HadError = resolver.ResolveStatements(statements)
	if HadError {
		return
	}
	hadRuntimeError = interpreter.Interpret(statements)
	if hadRuntimeError {
		return
	}

}
