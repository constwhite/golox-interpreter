package expr

import "github.com/constwhite/golox-interpreter/token"

type Expr struct{}
type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}
type Grouping struct {
	Expression Expr
}
type Literal struct {
	Value interface{}
}
type Unary struct {
	Operator token.Token
	Right    Expr
}
