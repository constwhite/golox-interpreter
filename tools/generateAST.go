package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var exprList = []string{
	"Binary:Left Expr, Operator token.Token, Right Expr",
	"Grouping:Expression Expr",
	"Literal:Value interface{}",
	"Unary:Operator token.Token, Right Expr",
}

func main() {
	//generates file with expression types
	args := os.Args
	if len(args) != 2 {
		log.Println("Usage: generateAST <output directory>")
	}
	outputDir := args[1]
	defineAST(outputDir, "expr", exprList)
}

func defineAST(outputDir string, baseName string, exprTypes []string) {
	path := fmt.Sprintf("%v/%v.go", outputDir, baseName)
	file, err := os.Create(path)
	if err != nil {
		log.Printf("error creating file: %v", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "package %v\n", baseName)
	fmt.Fprintf(file, `import "github.com/constwhite/golox-interpreter/token"`)
	fmt.Fprintln(file, "")
	baseNameUpper := fmt.Sprintf("%v%v", strings.ToUpper(string(baseName[0])), string(baseName[1:]))
	fmt.Fprintf(file, "type %v struct{}", baseNameUpper)

	for i := 0; i < len(exprTypes); i++ {
		exprType := exprTypes[i]
		structName := strings.Split(exprType, ":")[0]
		fields := strings.Split(exprType, ":")[1]
		defineType(file, structName, fields)
	}

}

func defineType(file *os.File, structName string, fieldList string) {
	fields := strings.Split(fieldList, ", ")
	fmt.Fprintf(file, "type %v struct {\n", structName)
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		fmt.Fprintf(file, "%v\n", field)
	}
	fmt.Fprintf(file, "}\n")
}
