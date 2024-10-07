package abstractSyntaxTree

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

type StmtVisitor interface {
	VisitExpressionStmt(stmt ExpressionStmt) interface{}
	VisitPrintStmt(stmt PrintStmt) interface{}
}
