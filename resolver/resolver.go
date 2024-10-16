package resolver

import (
	"io"

	"fmt"

	abs "github.com/constwhite/golox-interpreter/abstractSyntaxTree"
	"github.com/constwhite/golox-interpreter/errorHandler"

	in "github.com/constwhite/golox-interpreter/interpreter"
	t "github.com/constwhite/golox-interpreter/token"
)

type Resolver struct {
	interpreter    *in.Interpreter
	stdErr         io.Writer
	HadError       bool
	scopes         scopes
	currentFuntion functionType
	currentClass   classType
	ResolverError
}

type functionType uint8
type classType uint8

const (
	funcTypeNone functionType = iota
	funcTypeFunction
	funcTypeMethod
	funcTypeInitialiser
)

const (
	classTypeNone classType = iota
	classTypeClass
	classTypeSubclass
)

type ResolverError struct {
	error
	Line int
}

func (re *ResolverError) Error() error {
	return re.error
}
func NewResolver(interpreter *in.Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter}
}

//visit statements

func (r *Resolver) VisitBlockStmt(stmt abs.BlockStmt) interface{} {

	r.beginScope()
	r.ResolveStatements(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitVarStmt(stmt abs.VarStmt) interface{} {
	r.scopes.declare(stmt.Name)
	if stmt.Initialiser != nil {
		r.resolveExpr(stmt.Initialiser)
	}
	r.scopes.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt abs.FunctionStmt) interface{} {
	r.scopes.declare(stmt.Name)
	r.scopes.define(stmt.Name)
	r.resolveFunction(stmt, funcTypeFunction)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt abs.ExpressionStmt) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}
func (r *Resolver) VisitIfStmt(stmt abs.IfStmt) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}
func (r *Resolver) VisitPrintStmt(stmt abs.PrintStmt) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}
func (r *Resolver) VisitReturnStmt(stmt abs.ReturnStmt) interface{} {
	if r.currentFuntion == funcTypeNone {
		r.error(stmt.Keyword, "can not return from the top level code")
	}
	if stmt.Value != nil {
		if r.currentFuntion == funcTypeInitialiser {
			r.error(stmt.Keyword, "can not return from an initialiser")
		}
		r.resolveExpr(stmt.Value)
	}
	return nil
}
func (r *Resolver) VisitWhileStmt(stmt abs.WhileStmt) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitClassStmt(stmt abs.ClassStmt) interface{} {
	enclosingClass := r.currentClass
	r.currentClass = classTypeClass
	r.scopes.declare(stmt.Name)
	r.scopes.define(stmt.Name)
	if stmt.Superclass != nil && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		r.error(stmt.Superclass.Name, "a class can not inherit from itself")
		return nil
	}
	if stmt.Superclass != nil {
		r.currentClass = classTypeSubclass
		r.resolveExpr(stmt.Superclass)
	}
	if stmt.Superclass != nil {
		r.beginScope()
		defer r.endScope()
		r.scopes.peek()["super"] = true
	}

	r.beginScope()
	r.scopes.peek()["this"] = true
	for i := 0; i < len(stmt.Methods); i++ {
		method := stmt.Methods[i]
		declaration := funcTypeMethod
		if method.Name.Lexeme == "init" {
			declaration = funcTypeInitialiser
		}
		r.resolveFunction(method, declaration)
	}
	r.endScope()
	r.currentClass = enclosingClass
	return nil
}

// visit expressions
func (r *Resolver) VisitVariableExpr(expr abs.VariableExpr) interface{} {

	if len(r.scopes) > 0 {
		scope := r.scopes.peek()
		defined, declared := scope[expr.Name.Lexeme]
		if r.scopes.empty() && declared && !defined {
			r.error(expr.Name, "cant't read local variable in its own initialiser.")
		}

	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr abs.AssignExpr) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}
func (r *Resolver) VisitBinaryExpr(expr abs.BinaryExpr) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil
}
func (r *Resolver) VisitCallExpr(expr abs.CallExpr) interface{} {
	r.resolveExpr(expr.Callee)
	for i := 0; i < len(expr.Arguements); i++ {
		arg := expr.Arguements[i]
		r.resolveExpr(arg)
	}
	return nil
}
func (r *Resolver) VisitGroupingExpr(expr abs.GroupingExpr) interface{} {
	r.resolveExpr(expr.Expression)
	return nil
}
func (r *Resolver) VisitLiteralExpr(expr abs.LiteralExpr) interface{} {
	return nil
}
func (r *Resolver) VisitLogicalExpr(expr abs.LogicalExpr) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}
func (r *Resolver) VisitUnaryExpr(expr abs.UnaryExpr) interface{} {

	r.resolveExpr(expr.Right)
	return nil
}
func (r *Resolver) VisitGetExpr(expr abs.GetExpr) interface{} {
	r.resolveExpr(expr.Object)
	return nil
}
func (r *Resolver) VisitSetExpr(expr abs.SetExpr) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr abs.SuperExpr) interface{} {
	if r.currentClass == classTypeNone {
		r.error(expr.Keyword, "can't use 'super' outside of class")
	} else if r.currentClass != classTypeSubclass {
		r.error(expr.Keyword, "can't use 'super' in a class with no superclass")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitThisExpr(expr abs.ThisExpr) interface{} {
	if r.currentClass == classTypeNone {
		r.error(expr.Keyword, "can't use 'this' outside of a class")
		return nil
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

//helpers

func (r *Resolver) resolveFunction(function abs.FunctionStmt, fnType functionType) {
	enclosingFunction := r.currentFuntion
	r.currentFuntion = fnType

	r.beginScope()
	for i := 0; i < len(function.Params); i++ {
		param := function.Params[i]
		r.scopes.declare(param)
		r.scopes.define(param)
	}
	r.ResolveStatements(function.Body)
	r.endScope()
	r.currentFuntion = enclosingFunction
}

func (r *Resolver) resolveLocal(expr abs.Expr, name t.Token) {
	for i := r.scopes.size() - 1; i >= 0; i-- {
		scope := r.scopes[i]
		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.size()-1-i)
			return
		}
	}
}

// traverses list of statements and resolves the variables in each statement
func (r *Resolver) ResolveStatements(statements []abs.Stmt) bool {

	for i := 0; i < len(statements); i++ {
		stmt := statements[i]
		r.resolveStmt(stmt)
	}
	return r.HadError
}

// resolves a single statement
func (r *Resolver) resolveStmt(stmt abs.Stmt) {
	stmt.Accept(r)
}
func (r *Resolver) resolveExpr(expr abs.Expr) {
	expr.Accept(r)
}
func (r *Resolver) beginScope() {

	r.scopes.push(make(scope))
}

func (r *Resolver) endScope() {

	r.scopes.pop()
}

//errors

func (r *Resolver) error(token t.Token, msg string) {
	var where string
	if token.TokenType == t.TokenEOF {
		where = "at end"
	} else {
		where = fmt.Sprintf("at '%v'", token.Lexeme)
	}
	err := ResolverError{error: fmt.Errorf("%v %v", msg, where)}
	errorHandler.CompileError(r.stdErr, err.Error(), token.Line)
	r.ResolverError = err
	r.HadError = true
}
