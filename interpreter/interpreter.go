package interpreter

import (
	"errors"
	"fmt"
	"io"

	abs "github.com/constwhite/golox-interpreter/abstractSyntaxTree"
	env "github.com/constwhite/golox-interpreter/environment"
	e "github.com/constwhite/golox-interpreter/errorHandler"
	t "github.com/constwhite/golox-interpreter/token"
)

type Interpreter struct {
	stdErr          io.Writer
	stdOut          io.Writer
	HasRuntimeError bool
	RuntimeError    runtimeError
	Environment     *env.Environment
	Globals         *env.Environment
	Locals          map[abs.Expr]int
}

type runtimeError struct {
	error
	Line int
}

func (rte *runtimeError) Error() error {
	return rte.error
}

func NewInterpreter(stdErr io.Writer, stdOut io.Writer) *Interpreter {
	global := env.NewEnvironment(nil)
	global.Define("clock", Clock{})
	return &Interpreter{stdErr: stdErr, stdOut: stdOut, Environment: global, Globals: global, Locals: make(map[abs.Expr]int)}
}

func (i *Interpreter) Interpret(stmtList []abs.Stmt) (HasRuntimeError bool) {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(runtimeError); ok {
				HasRuntimeError = true
				return
			} else {
				panic(err)
			}
		}
	}()
	for index := 0; index < len(stmtList); index++ {
		stmt := stmtList[index]
		i.execute(stmt)
	}
	return HasRuntimeError

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
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)

	case t.TokenGreaterEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)

	case t.TokenLesser:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)

	case t.TokenLesserEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)

	case t.TokenBangEqual:
		return left != right
	case t.TokenEqualEqual:
		return left == right
	case t.TokenMinus:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)

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
		e.RuntimeError(i.stdErr, err.Error(), err.Line)
		panic(err)
	case t.TokenSlash:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)

	case t.TokenStar:
		i.checkNumberOperands(expr.Operator, left, right)

		return left.(float64) * right.(float64)

	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr abs.VariableExpr) interface{} {
	// value, err := i.Environment.Get(expr.Name)
	value, err := i.lookupVariable(expr.Name, expr)
	if err != nil {
		runtimeErr := runtimeError{error: err, Line: expr.Name.Line}
		i.RuntimeError = runtimeErr

	}
	return value
}
func (i *Interpreter) VisitAssignExpr(expr abs.AssignExpr) interface{} {
	value := i.evaluate(expr.Value)
	distance, ok := i.Locals[expr]
	if ok {
		i.Environment.AssignAt(distance, expr.Name, value)
	} else {
		if err := i.Globals.Assign(expr.Name, value); err != nil {
			runtimeErr := runtimeError{error: err, Line: expr.Name.Line}
			i.RuntimeError = runtimeErr
		}
	}
	return value
}

func (i *Interpreter) VisitLogicalExpr(expr abs.LogicalExpr) interface{} {
	left := i.evaluate(expr.Left)
	if expr.Operator.TokenType == t.TokenOr {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr abs.CallExpr) interface{} {
	callee := i.evaluate(expr.Callee)
	var arguements []interface{}
	for index := 0; index < len(expr.Arguements); index++ {
		arguement := expr.Arguements[index]
		arguements = append(arguements, i.evaluate(arguement))
	}
	function, callable := callee.(loxCallable)
	if !callable {
		err := runtimeError{error: errors.New("can only call funtions and classes"), Line: expr.Paren.Line}
		i.RuntimeError = err
		return nil
	}
	if len(arguements) != function.arity() {
		err := runtimeError{error: fmt.Errorf("expected %v arguements but got %v", function.arity(), len(arguements)), Line: expr.Paren.Line}
		i.RuntimeError = err
		return nil
	}
	return function.call(i, arguements)
}

func (i *Interpreter) VisitGetExpr(expr abs.GetExpr) interface{} {
	object := i.evaluate(expr.Object)
	instance, isInstance := object.(loxInstance)
	if !isInstance {
		err := runtimeError{error: errors.New("only instances have properties"), Line: expr.Name.Line}
		i.RuntimeError = err
		return nil
	}
	property, err := instance.get(expr.Name)
	if err != nil {
		i.RuntimeError = runtimeError{error: err, Line: expr.Name.Line}
		return nil
	}
	return property
}

func (i *Interpreter) VisitSetExpr(expr abs.SetExpr) interface{} {
	object := i.evaluate(expr.Object)
	instance, isInstance := object.(loxInstance)
	if !isInstance {
		err := runtimeError{error: errors.New("only instances have fields"), Line: expr.Name.Line}
		i.RuntimeError = err
		return nil
	}
	value := i.evaluate(expr.Value)
	instance.set(expr.Name, value)
	return value
}

func (i *Interpreter) VisitSuperExpr(expr abs.SuperExpr) interface{} {
	distance := i.Locals[expr]
	superclass := i.Environment.GetAt(distance, "super").(*loxClass)
	object := i.Environment.GetAt(distance-1, "this").(*loxInstance)
	method := superclass.findMethod(expr.Method.Lexeme)
	if method == nil {
		err := runtimeError{error: fmt.Errorf("undefined property %v", expr.Method.Lexeme), Line: expr.Method.Line}
		i.RuntimeError = err
		return nil
	}
	return method.bind(object)
}

func (i *Interpreter) VisitThisExpr(expr abs.ThisExpr) interface{} {
	value, err := i.lookupVariable(expr.Keyword, expr)
	if err != nil {
		runtimeErr := runtimeError{error: err, Line: expr.Keyword.Line}
		i.RuntimeError = runtimeErr

	}
	return value
}

// statement visitors
func (i *Interpreter) VisitExpressionStmt(stmt abs.ExpressionStmt) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt abs.FunctionStmt) interface{} {
	function := loxFunction{Declaration: stmt, Closure: i.Environment, isInitialiser: false}
	i.Environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt abs.PrintStmt) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Fprint(i.stdOut, i.stringify(value))
	return nil
}

type returnValue struct {
	Value interface{}
}

func (i *Interpreter) VisitReturnStmt(stmt abs.ReturnStmt) interface{} {
	var value interface{} = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	returnValue := returnValue{Value: value}
	panic(returnValue)
}

func (i *Interpreter) VisitVarStmt(stmt abs.VarStmt) interface{} {
	var value interface{}
	if stmt.Initialiser != nil {
		value = i.evaluate(stmt.Initialiser)
	}
	i.Environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt abs.BlockStmt) interface{} {
	i.executeBlock(stmt.Statements, env.NewEnvironment(i.Environment))
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt abs.IfStmt) interface{} {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt abs.WhileStmt) interface{} {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt abs.ClassStmt) interface{} {
	var superclass *loxClass = nil
	if stmt.Superclass != nil {
		superclassInterface := i.evaluate(stmt.Superclass)
		superclassAssert, ok := superclassInterface.(loxClass)
		if !ok {
			err := runtimeError{error: errors.New("superclass must be a class"), Line: stmt.Name.Line}
			i.RuntimeError = err
			return nil
		}
		superclass = &superclassAssert
	}

	i.Environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.Environment = env.NewEnvironment(i.Environment)
		i.Environment.Define("super", superclass)
	}

	methods := make(map[string]loxFunction)
	for index := 0; index < len(stmt.Methods); index++ {
		method := stmt.Methods[index]
		isInit := method.Name.Lexeme == "init"
		function := loxFunction{Declaration: method, Closure: i.Environment, isInitialiser: isInit}
		methods[method.Name.Lexeme] = function
	}

	class := loxClass{Name: stmt.Name.Lexeme, SuperClass: superclass, methods: methods}

	if superclass != nil {
		i.Environment = i.Environment.Enclosing
	}

	i.Environment.Assign(stmt.Name, class)
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
	return fmt.Sprint(value)
}

func (i *Interpreter) lookupVariable(name t.Token, expr abs.Expr) (interface{}, error) {
	distance, ok := i.Locals[expr]
	if ok {
		return i.Environment.GetAt(distance, name.Lexeme), nil
	} else {
		return i.Globals.Get(name)
	}
}

func (i *Interpreter) execute(stmt abs.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) Resolve(expr abs.Expr, depth int) {
	i.Locals[expr] = depth
}

func (i *Interpreter) executeBlock(statements []abs.Stmt, environment *env.Environment) {
	previous := i.Environment
	defer func() {
		i.Environment = previous
	}()
	i.Environment = environment
	for index := 0; index < len(statements); index++ {
		stmt := statements[index]
		i.execute(stmt)
	}
}

// runtime errors

func (i *Interpreter) checkNumberOperand(operator t.Token, operand interface{}) bool {
	if _, ok := operand.(float64); ok {
		return true
	}
	err := runtimeError{error: fmt.Errorf("operand must be a number"), Line: operator.Line}

	i.RuntimeError = err
	e.RuntimeError(i.stdErr, err.Error(), err.Line)
	panic(err)
}

func (i *Interpreter) checkNumberOperands(operator t.Token, left interface{}, right interface{}) bool {
	_, leftIsFloat := left.(float64)
	_, rightIsFloat := right.(float64)
	if leftIsFloat && rightIsFloat {
		return true
	}
	err := runtimeError{error: fmt.Errorf("operands must be numbers"), Line: operator.Line}
	e.RuntimeError(i.stdErr, err.Error(), err.Line)
	i.RuntimeError = err

	panic(err)

}
