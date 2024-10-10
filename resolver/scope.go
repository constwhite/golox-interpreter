package resolver

import t "github.com/constwhite/golox-interpreter/token"

type scope map[string]bool

// scopes is implemented as a stack built on top of an array
type scopes []scope

func (s *scopes) declare(token t.Token) {
	if s.empty() {
		return
	}
	scope := s.peek()
	scope[token.Lexeme] = false
}

func (s *scopes) define(token t.Token) {
	if s.empty() {
		return
	}
	scope := s.peek()
	scope[token.Lexeme] = true
}

// stack methods. scopes is recieved as a pointer as it it modifying the original scopes slice on the resolver
func (s *scopes) push(scope scope) {
	*s = append(*s, scope)
}

func (s *scopes) pop() scope {
	if s.empty() {
		return nil
	}
	pop := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return pop
}

func (s *scopes) peek() scope {
	return (*s)[len(*s)-1]
}

func (s *scopes) empty() bool {
	if len(*s) == 0 {
		return true
	} else {
		return false
	}
}
func (s *scopes) size() int {
	return len(*s)
}
