package parser

import (
	"fmt"
	"io"

	abs "github.com/constwhite/golox-interpreter/abstractSyntaxTree"
	e "github.com/constwhite/golox-interpreter/errorHandler"
	t "github.com/constwhite/golox-interpreter/token"
)

type Parser struct {
	stdErr       io.Writer
	current      int
	sourceTokens []t.Token
	HadError     bool
	parseError
}

type parseError struct {
	msg string
}

func (pe *parseError) Error() string {
	return pe.msg
}

func NewParser(sourceTokens []t.Token, stdErr io.Writer) *Parser {
	return &Parser{sourceTokens: sourceTokens, stdErr: stdErr}
}

func (p *Parser) Parse() ([]abs.Stmt, bool) {
	var statements []abs.Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	if p.Error() != "" {
		p.HadError = true
		return nil, p.HadError
	}
	return statements, p.HadError
}

// grammar functions
func (p *Parser) declaration() abs.Stmt {
	if p.Error() != "" {
		p.HadError = true
		return nil
	}
	if p.match(t.TokenVar) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() abs.Stmt {
	name := p.consume(t.TokenIdentifier, "expect variable name")
	var initialiser abs.Expr = nil
	if p.match(t.TokenEqual) {
		initialiser = p.expression()
	}
	p.consume(t.TokenSemiColon, "expect ';' after variable declaration")
	return abs.VarStmt{Name: name, Initialiser: initialiser}
}

func (p *Parser) statement() abs.Stmt {
	if p.match(t.TokenPrint) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() abs.Stmt {
	value := p.expression()
	p.consume(t.TokenSemiColon, "expect ';' after value")
	return abs.PrintStmt{Expression: value}
}
func (p *Parser) expressionStatement() abs.Stmt {
	expression := p.expression()
	p.consume(t.TokenSemiColon, "expect ';' after value")
	return abs.ExpressionStmt{Expression: expression}
}

func (p *Parser) expression() abs.Expr {
	return p.equality()
}
func (p *Parser) equality() abs.Expr {
	expr := p.comparison()
	for p.match(t.TokenBangEqual, t.TokenEqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = abs.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}
func (p *Parser) comparison() abs.Expr {
	expr := p.term()

	for p.match(t.TokenGreater, t.TokenGreaterEqual, t.TokenLesser, t.TokenLesserEqual) {
		operator := p.previous()
		right := p.term()
		expr = abs.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() abs.Expr {
	expr := p.factor()
	for p.match(t.TokenMinus, t.TokenPlus) {
		operator := p.previous()
		right := p.factor()
		expr = abs.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() abs.Expr {
	expr := p.unary()

	for p.match(t.TokenSlash, t.TokenStar) {
		operator := p.previous()
		right := p.unary()
		expr = abs.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() abs.Expr {
	if p.match(t.TokenBang, t.TokenMinus) {
		operator := p.previous()
		right := p.unary()
		return abs.UnaryExpr{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() abs.Expr {
	if p.match(t.TokenFalse) {
		return abs.LiteralExpr{Value: false}
	}
	if p.match(t.TokenTrue) {
		return abs.LiteralExpr{Value: true}
	}
	if p.match(t.TokenNil) {
		return abs.LiteralExpr{Value: nil}
	}

	if p.match(t.TokenNumber, t.TokenString) {
		return abs.LiteralExpr{Value: p.previous().Literal}
	}
	if p.match(t.TokenIdentifier) {
		return abs.VariableExpr{Name: p.previous()}
	}

	if p.match(t.TokenLeftParen) {
		expr := p.expression()
		p.consume(t.TokenRightParen, "expect ')' after expression")
		return abs.GroupingExpr{Expression: expr}
	}
	p.error(p.peek(), "expect expression")
	return nil
}

// error handling
func (p *Parser) consume(tokenType t.TokenType, message string) t.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return t.Token{}

}

func (p *Parser) error(token t.Token, message string) {
	var where string
	if token.TokenType == t.TokenEOF {
		where = "at end"
	} else {
		where = fmt.Sprintf("at '%v'", token.Lexeme)
	}
	err := parseError{msg: message}
	e.ReportError(p.stdErr, err.Error(), where, token.Line+1)
	p.parseError = err
	// panic(err)
}

func (p *Parser) synchronise() {
	for !p.isAtEnd() {
		if p.previous().TokenType == t.TokenSemiColon {
			return
		}
		switch p.peek().TokenType {
		case t.TokenClass:
		case t.TokenFun:
		case t.TokenVar:
		case t.TokenFor:
		case t.TokenIf:
		case t.TokenWhile:
		case t.TokenPrint:
		case t.TokenReturn:
			return

		}
	}

	p.advance()
}

// helper functions
func (p *Parser) match(tokenTypes ...t.TokenType) bool {
	for i := 0; i < len(tokenTypes); i++ {
		tokenType := tokenTypes[i]
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType t.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() t.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == t.TokenEOF
}
func (p *Parser) peek() t.Token {
	return p.sourceTokens[p.current]
}
func (p *Parser) previous() t.Token {
	return p.sourceTokens[p.current-1]
}
