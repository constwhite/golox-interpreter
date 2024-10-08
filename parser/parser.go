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

func (p *Parser) whileStatement() abs.Stmt {
	p.consume(t.TokenLeftParen, "expect '(' after 'while'")
	condition := p.expression()
	p.consume(t.TokenRightParen, "expect ')' after condition")
	body := p.statement()
	return abs.WhileStmt{Condition: condition, Body: body}
}

func (p *Parser) statement() abs.Stmt {
	if p.match(t.TokenFor) {
		return p.forStatement()
	}
	if p.match(t.TokenIf) {
		return p.ifStatement()
	}
	if p.match(t.TokenPrint) {
		return p.printStatement()
	}
	if p.match(t.TokenWhile) {
		return p.whileStatement()
	}
	if p.match(t.TokenLeftBrace) {
		return abs.BlockStmt{Statements: p.blockStatement()}
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() abs.Stmt {
	p.consume(t.TokenLeftParen, "expect '(' after 'for'")
	var initialiser abs.Stmt
	if p.match(t.TokenSemiColon) {
		initialiser = nil
	} else if p.match(t.TokenVar) {
		initialiser = p.varDeclaration()
	} else {
		initialiser = p.expressionStatement()
	}
	var condition abs.Expr = nil
	if !p.check(t.TokenSemiColon) {
		condition = p.expression()
	}
	p.consume(t.TokenSemiColon, "expect ';' after loop condition")
	var increment abs.Expr = nil
	if !p.check(t.TokenRightParen) {
		increment = p.expression()
	}
	p.consume(t.TokenRightParen, "expect ')' after for clauses")
	body := p.statement()

	if increment != nil {
		body = abs.BlockStmt{
			Statements: []abs.Stmt{body, abs.ExpressionStmt{Expression: increment}},
		}
	}
	if condition == nil {
		condition = abs.LiteralExpr{Value: true}
		body = abs.WhileStmt{Condition: condition, Body: body}
	}
	if initialiser != nil {
		body = abs.BlockStmt{Statements: []abs.Stmt{initialiser, body}}
	}

	return body
}

func (p *Parser) ifStatement() abs.Stmt {
	p.consume(t.TokenLeftParen, "expect '(' after 'if'")
	condition := p.expression()
	p.consume(t.TokenRightParen, "expect ')' after if condition")
	thenBranch := p.statement()
	var elseBranch abs.Stmt = nil
	if p.match(t.TokenElse) {
		elseBranch = p.statement()
	}
	return abs.IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}

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
func (p *Parser) blockStatement() []abs.Stmt {
	var statements []abs.Stmt = nil
	for !p.check(t.TokenRightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(t.TokenRightBrace, "expect '}' after block")
	return statements
}

func (p *Parser) expression() abs.Expr {
	return p.assignment()
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

	return p.call()
}

func (p *Parser) call() abs.Expr {
	expr := p.primary()
	for {
		if p.match(t.TokenLeftParen) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee abs.Expr) abs.Expr {
	var arguements []abs.Expr = nil
	if !p.check(t.TokenRightParen) {
		for {
			if len(arguements) >= 255 {
				p.error(p.peek(), "functions can not accept more than 255 arguements")
			}
			arguements = append(arguements, p.expression())
			if !p.match(t.TokenComma) {
				break
			}
		}
	}
	paren := p.consume(t.TokenRightParen, "expect ')' after arguements")
	return abs.CallExpr{Callee: callee, Paren: paren, Arguements: arguements}
}

func (p *Parser) assignment() abs.Expr {
	expr := p.or()
	if p.match(t.TokenEqual) {
		equals := p.previous()
		value := p.assignment()
		if exprVariable, ok := expr.(abs.VariableExpr); ok {
			name := exprVariable.Name
			return abs.AssignExpr{Name: name, Value: value}
		}

		p.error(equals, "invalid assignment target")
	}
	return expr
}

func (p *Parser) or() abs.Expr {
	expr := p.and()
	for p.match(t.TokenOr) {
		operator := p.previous()
		right := p.and()
		expr = abs.LogicalExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) and() abs.Expr {
	expr := p.equality()
	for p.match(t.TokenAnd) {
		operator := p.previous()
		right := p.equality()
		expr = abs.LogicalExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
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
