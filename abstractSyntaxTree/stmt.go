package abstractSyntaxTree

import "github.com/constwhite/golox-interpreter/token"

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
	Name        token.Token
}

func (s VarStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(s)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt ExpressionStmt) interface{}
	VisitPrintStmt(stmt PrintStmt) interface{}
	VisitVarStmt(stmt VarStmt) interface{}
}
