package abstractSyntaxTree

import "fmt"

type Printer struct {
	ExprVisitor
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *Printer) parenthesise(name string, exprs ...Expr) string {
	stringBuilder := fmt.Sprintf("(%v", name)
	for i := 0; i < len(exprs); i++ {
		expr := exprs[i]
		stringBuilder = fmt.Sprintf("%v %v", stringBuilder, expr.Accept(p))
	}
	stringBuilder = fmt.Sprintf("%v)", stringBuilder)
	return stringBuilder

}
func (p *Printer) VisitBinaryExpr(expr BinaryExpr) interface{} {
	return p.parenthesise(expr.Operator.Lexeme, expr.Left, expr.Right)
}
func (p *Printer) VisitGroupingExpr(expr GroupingExpr) interface{} {
	return p.parenthesise("group", expr.Expression)
}
func (p *Printer) VisitLiteralExpr(expr LiteralExpr) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Value)
}
func (p *Printer) VisitUnaryExpr(expr UnaryExpr) interface{} {
	return p.parenthesise(expr.Operator.Lexeme, expr.Right)
}
