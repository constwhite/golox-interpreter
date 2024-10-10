package abstractSyntaxTree

import t "github.com/constwhite/golox-interpreter/token"

type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}
type BinaryExpr struct {
	Left     Expr
	Operator t.Token
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
	Operator t.Token
	Right    Expr
}

func (e UnaryExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(e)
}

type VariableExpr struct {
	Name t.Token
}

func (e VariableExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(e)
}

type AssignExpr struct {
	Name  t.Token
	Value Expr
}

func (e AssignExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(e)
}

type LogicalExpr struct {
	Left     Expr
	Operator t.Token
	Right    Expr
}

func (e LogicalExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicalExpr(e)
}

type CallExpr struct {
	Callee     Expr
	Paren      t.Token
	Arguements []Expr
}

func (e CallExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCallExpr(e)
}

type GetExpr struct {
	Object Expr
	Name   t.Token
}

func (e GetExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGetExpr(e)
}

type SetExpr struct {
	Object Expr
	Name   t.Token
	Value  Expr
}

func (e SetExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitSetExpr(e)
}

type ThisExpr struct {
	Keyword t.Token
}

func (e ThisExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitThisExpr(e)
}

type ExprVisitor interface {
	VisitBinaryExpr(expr BinaryExpr) interface{}
	VisitGroupingExpr(expr GroupingExpr) interface{}
	VisitLiteralExpr(expr LiteralExpr) interface{}
	VisitUnaryExpr(expr UnaryExpr) interface{}
	VisitVariableExpr(expr VariableExpr) interface{}
	VisitAssignExpr(expr AssignExpr) interface{}
	VisitLogicalExpr(expr LogicalExpr) interface{}
	VisitCallExpr(expr CallExpr) interface{}
	VisitGetExpr(expr GetExpr) interface{}
	VisitSetExpr(expr SetExpr) interface{}
	VisitThisExpr(expr ThisExpr) interface{}
}
