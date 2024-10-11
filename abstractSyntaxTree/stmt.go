package abstractSyntaxTree

import (
	t "github.com/constwhite/golox-interpreter/token"
)

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}
type ExpressionStmt struct {
	Expression Expr
}

func (s ExpressionStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(s)
}

type PrintStmt struct {
	Expression Expr
}

func (s PrintStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(s)
}

type VarStmt struct {
	Initialiser Expr
	Name        t.Token
}

func (s VarStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(s)
}

type BlockStmt struct {
	Statements []Stmt
}

func (s BlockStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlockStmt(s)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (s IfStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitIfStmt(s)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (s WhileStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitWhileStmt(s)
}

type FunctionStmt struct {
	Name   t.Token
	Params []t.Token
	Body   []Stmt
}

func (s FunctionStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitFunctionStmt(s)
}

type ReturnStmt struct {
	Keyword t.Token
	Value   Expr
}

func (s ReturnStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitReturnStmt(s)
}

type ClassStmt struct {
	Name       t.Token
	Superclass *VariableExpr
	Methods    []FunctionStmt
}

func (s ClassStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitClassStmt(s)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt ExpressionStmt) interface{}
	VisitPrintStmt(stmt PrintStmt) interface{}
	VisitVarStmt(stmt VarStmt) interface{}
	VisitBlockStmt(stmt BlockStmt) interface{}
	VisitIfStmt(stmt IfStmt) interface{}
	VisitWhileStmt(stmt WhileStmt) interface{}
	VisitFunctionStmt(stmt FunctionStmt) interface{}
	VisitReturnStmt(stmt ReturnStmt) interface{}
	VisitClassStmt(stmt ClassStmt) interface{}
}
