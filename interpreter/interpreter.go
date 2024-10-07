package interpreter

import (
	"fmt"
	"io"

	abs "github.com/constwhite/golox-interpreter/abstractSyntaxTree"
	e "github.com/constwhite/golox-interpreter/errorHandler"
	t "github.com/constwhite/golox-interpreter/token"
)

type Interpreter struct {
	stdErr          io.Writer
	stdOut          io.Writer
	HasRuntimeError bool
	RuntimeError    runtimeError
}

type runtimeError struct {
	error
	Line int
}

func (rte *runtimeError) Error() error {
	return rte.error
}

func NewInterpreter(stdErr io.Writer, stdOut io.Writer) *Interpreter {
	return &Interpreter{stdErr: stdErr, stdOut: stdOut}
}

func (i *Interpreter) Interpret(stmtList []abs.Stmt) bool {
	if err := i.RuntimeError.Error(); err != nil {
		e.RuntimeError(i.stdErr, err, i.RuntimeError.Line)
		i.HasRuntimeError = true
		return i.HasRuntimeError
	} else {
		for index := 0; index < len(stmtList); index++ {
			stmt := stmtList[index]
			i.execute(stmt)
		}
		return false
	}
}

func (i *Interpreter) execute(stmt abs.Stmt) {
	stmt.Accept(i)
}

//expression visitors

func (i *Interpreter) VisitLiteralExpr(expr abs.LiteralExpr) interface{} {
	return expr.Value
}
func (i *Interpreter) VisitGroupingExpr(expr abs.GroupingExpr) interface{} {
	return i.evaluate(expr.Expression)
}
func (i *Interpreter) VisitUnaryExpr(expr abs.UnaryExpr) interface{} {
	right := i.evaluate(expr.Right)
	switch expr.Operator.TokenType {
	case t.TokenBang:
		return !i.isTruthy(right)
	case t.TokenMinus:
		if i.checkNumberOperand(expr.Operator, right) {
			return -right.(float64)

		}
	}
	return nil
}
func (i *Interpreter) VisitBinaryExpr(expr abs.BinaryExpr) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case t.TokenGreater:
		if i.checkNumberOperands(expr.Operator, left, right) {
			return left.(float64) > right.(float64)
		}
	case t.TokenGreaterEqual:
		if i.checkNumberOperands(expr.Operator, left, right) {
			return left.(float64) >= right.(float64)

		}
	case t.TokenLesser:
		if i.checkNumberOperands(expr.Operator, left, right) {
			return left.(float64) < right.(float64)

		}
	case t.TokenLesserEqual:
		if i.checkNumberOperands(expr.Operator, left, right) {

			return left.(float64) <= right.(float64)
		}
	case t.TokenBangEqual:
		return left != right
	case t.TokenEqualEqual:
		return left == right
	case t.TokenMinus:
		if i.checkNumberOperands(expr.Operator, left, right) {
			return left.(float64) - right.(float64)

		}
	case t.TokenPlus:
		_, leftIsFloat := left.(float64)
		_, rightIsFloat := right.(float64)
		if leftIsFloat && rightIsFloat {
			return left.(float64) + right.(float64)
		}
		_, leftIsString := left.(string)
		_, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return left.(string) + right.(string)
		}
		err := runtimeError{error: fmt.Errorf("operands must be numbers or string"), Line: expr.Operator.Line}
		i.RuntimeError = err
	case t.TokenSlash:
		if i.checkNumberOperands(expr.Operator, left, right) {
			return left.(float64) / right.(float64)

		}
	case t.TokenStar:
		if i.checkNumberOperands(expr.Operator, left, right) {

			return left.(float64) * right.(float64)
		}
	}

	return nil
}

// statement visitors
func (i *Interpreter) VisitExpressionStmt(stmt abs.ExpressionStmt) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt abs.PrintStmt) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Fprint(i.stdOut, i.stringify(value))
	return nil
}

// helpers
func (i *Interpreter) evaluate(expr abs.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	} // nil value returns false
	if v, ok := object.(bool); ok { //if false return v. any other value returns !ok
		return v
	}
	return true //everything else is truthy
}

func (i *Interpreter) stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}
	// if _, isNumber := value.(float64); isNumber {
	// 	text, _ := value.(string)
	// 	// if text[len(text)-2:] == ".0" {
	// 	// 	text = text[:len(text)-2]
	// 	// }
	// 	return text
	// }
	// return value.(string)
	return fmt.Sprint(value)
}

// runtime errors
func (i *Interpreter) checkNumberOperand(operator t.Token, operand interface{}) bool {
	if _, ok := operand.(float64); ok {
		return true
	}
	err := runtimeError{error: fmt.Errorf("operand must be a number"), Line: operator.Line}
	i.RuntimeError = err

	return false
}

func (i *Interpreter) checkNumberOperands(operator t.Token, left interface{}, right interface{}) bool {
	_, leftIsFloat := left.(float64)
	_, rightIsFloat := right.(float64)
	if leftIsFloat && rightIsFloat {
		return true
	}
	err := runtimeError{error: fmt.Errorf("operands must be numbers"), Line: operator.Line}
	i.RuntimeError = err

	return false
}
