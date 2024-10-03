package main

import (
	"fmt"

	abstractsyntaxtree "github.com/constwhite/golox-interpreter/abstractSyntaxTree"
	"github.com/constwhite/golox-interpreter/token"
)

func main() {
	expression := abstractsyntaxtree.BinaryExpr{
		Left: abstractsyntaxtree.UnaryExpr{
			Operator: token.Token{TokenType: token.TokenMinus, Lexeme: "-", Literal: nil, Line: 1},
			Right:    abstractsyntaxtree.LiteralExpr{Value: 123},
		},
		Operator: token.Token{TokenType: token.TokenStar, Lexeme: "*", Literal: nil, Line: 1},
		Right: abstractsyntaxtree.GroupingExpr{
			Expression: abstractsyntaxtree.LiteralExpr{Value: 45.67},
		},
	}
	printer := abstractsyntaxtree.NewPrinter()
	fmt.Println(printer.Print(expression))
}
