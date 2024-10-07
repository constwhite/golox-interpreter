package abstractSyntaxTree

import "github.com/constwhite/golox-interpreter/token"

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}
type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e BinaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e GroupingExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e LiteralExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (e UnaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(e)
}

type VariableExpr struct {
	Name token.Token
}

func (e VariableExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(e)
}

type ExprVisitor interface {
	VisitBinaryExpr(expr BinaryExpr) interface{}
	VisitGroupingExpr(expr GroupingExpr) interface{}
	VisitLiteralExpr(expr LiteralExpr) interface{}
	VisitUnaryExpr(expr UnaryExpr) interface{}
	VisitVariableExpr(expr VariableExpr) interface{}
}
