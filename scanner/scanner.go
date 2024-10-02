package scanner

import (
	"fmt"
	"io"
	"strconv"

	"github.com/constwhite/golox-interpreter/errorHandler"
	"github.com/constwhite/golox-interpreter/token"
)

type Scanner struct {
	//source code stored on scanner struct as string
	source  string
	tokens  []token.Token
	current int
	start   int
	line    int
	stdErr  io.Writer
}

func NewScanner(source string, stdErr io.Writer) *Scanner {
	return &Scanner{source: source, stdErr: stdErr}
}

var keywords = map[string]token.TokenType{
	"and":    token.TokenAnd,
	"class":  token.TokenClass,
	"else":   token.TokenElse,
	"false":  token.TokenFalse,
	"for":    token.TokenFor,
	"fun":    token.TokenFun,
	"if":     token.TokenIf,
	"nil":    token.TokenNil,
	"or":     token.TokenOr,
	"print":  token.TokenPrint,
	"return": token.TokenReturn,
	"super":  token.TokenSuper,
	"this":   token.TokenThis,
	"true":   token.TokenTrue,
	"var":    token.TokenVar,
	"while":  token.TokenWhile,
}

func (s *Scanner) ScanTokens() []token.Token {
	//loop through source until reaching the end then appends one End of file (EOF) token
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.Token{TokenType: token.TokenEOF, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	//single character lexemes
	case '(':
		s.addToken(token.TokenLeftParen)
	case ')':
		s.addToken(token.TokenRightParen)
	case '{':
		s.addToken(token.TokenLeftBrace)
	case '}':
		s.addToken(token.TokenRightBrace)
	case ',':
		s.addToken(token.TokenComma)
	case '.':
		s.addToken(token.TokenDot)
	case '-':
		s.addToken(token.TokenMinus)
	case '+':
		s.addToken(token.TokenPlus)
	case ';':
		s.addToken(token.TokenSemiColon)
	case '*':
		s.addToken(token.TokenStar)
	//operators
	case '!':

		if s.match('=') {
			s.addToken(token.TokenBangEqual)
		} else {
			s.addToken(token.TokenBang)
		}
	case '=':

		if s.match('=') {
			s.addToken(token.TokenEqualEqual)
		} else {
			s.addToken(token.TokenEqual)
		}
	case '>':

		if s.match('=') {
			s.addToken(token.TokenGreaterEqual)
		} else {
			s.addToken(token.TokenGreater)
		}
	case '<':

		if s.match('=') {
			s.addToken(token.TokenLesserEqual)
		} else {
			s.addToken(token.TokenLesser)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.TokenSlash)
		}
	//ignore whitespace
	case ' ':
	case '\r':
	case '\t':
	//line break
	case '\n':
		s.line++
	//strings
	case '"':
		s.string()
	//defaults
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			errorHandler.ReportError(s.stdErr, fmt.Sprintf("Unexpected Character: %v", string(c)), s.line)
			// s.error(fmt.Sprintf("Unexpected Character: %v", string(c)))
		}
	}
}

func (s *Scanner) advance() rune {
	//returns the rune of the current byte in the source string and moves to the next byte
	char := rune(s.source[s.current])
	s.current++
	return char
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		//iterates over bytes in source segment until it reaches a closing ' " ' or reaches the end of the lexeme.
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		//if no closing ' " ' return error
		errorHandler.ReportError(s.stdErr, "Unterminated string", s.line)
		// s.error("Unterminated string")
		return
	}

	s.advance()
	// add token
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(token.TokenString, value)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		s.addToken(token.TokenIdentifier)

	} else {

		s.addToken(tokenType)
	}
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isDigit(c) || s.isAlpha(c)
}

func (s *Scanner) number() {
	//loop until non digit character
	for s.isDigit(s.peek()) {
		s.advance()
	}
	//check if first non digit is '.'. if true, peek char after '.' if digit, keep looping
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	//if '.' not followed by digit or end of number add token
	//convert into float64
	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		panic(err)
	}
	s.addTokenWithLiteral(token.TokenNumber, value)

}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if expected != rune(s.source[s.current]) {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(
		s.tokens,
		token.Token{
			TokenType: tokenType,
			Lexeme:    text,
			Literal:   literal,
			Line:      s.line,
		})

}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)

}

// func (s *Scanner) error(msg string) {
// 	// _, _ = s.stdErr.Write([]byte(fmt.Sprintf("[line: %v] Error: %s\n", s.line, msg)))
// 	fmt.Fprintf(s.stdErr, "[line: %v] Error: %s\n", s.line, msg)
// }
