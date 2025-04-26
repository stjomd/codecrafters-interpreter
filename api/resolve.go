package api

import (
	intp "github.com/codecrafters-io/interpreter-starter-go/api/interpreter"
	"github.com/codecrafters-io/interpreter-starter-go/spec"
)

type resolver struct { // implements spec.ExprVisitor[any, error], spec.StmtVisitor[error]
	intp *intp.Interpreter
	scopes stack[map[string]bool]
	hadError bool
	currentFuncType intp.FunctionType
	currentClassType intp.ClassType
}

type stack[T any] struct {
	slice []T
}
func (s *stack[T]) push(elem T) {
	s.slice = append(s.slice, elem)
}
func (s *stack[T]) peek() T {
	if s.size() > 0 {
		index := len(s.slice) - 1
		elem := s.slice[index]
		return elem
	}
	var zero T
	return zero
}
func (s *stack[T]) pop() T {
	if s.size() > 0 {
		index := len(s.slice) - 1
		elem := s.slice[index]
		s.slice = s.slice[:index]
		return elem
	}
	var zero T
	return zero
}
func (s *stack[T]) get(index int) T {
	return s.slice[index]
}
func (s *stack[T]) size() int {
	return len(s.slice)
}
func (s *stack[T]) isEmpty() bool {
	return s.size() == 0
}

// MARK: - Methods

func (rslv *resolver) resolveStmts(stmts *[]spec.Stmt) {
	for _, stmt := range *stmts {
		rslv.resolveStmt(stmt)
	}
}

func (rslv *resolver) resolveStmt(stmt spec.Stmt) {
	stmt.Exec(rslv)
}

func (rslv *resolver) resolveExpr(expr spec.Expr) {
	expr.Eval(rslv)
}

func (rslv *resolver) resolveFunction(fs spec.FuncStmt, funcType intp.FunctionType) {
	origFuncType := rslv.currentFuncType
	rslv.currentFuncType = funcType
	defer func() { rslv.currentFuncType = origFuncType }()

	rslv.beginScope()
	for _, param := range fs.Params {
		rslv.declare(param);
    rslv.define(param);
	}
	rslv.resolveStmts(&fs.Body);
	rslv.endScope()
}

func (rslv *resolver) resolveLocal(expr spec.Expr, name spec.Token) {
	for i := rslv.scopes.size() - 1; i >= 0; i-- {
		if _, contains := rslv.scopes.get(i)[name.Lexeme]; contains {
			rslv.intp.Resolve(expr, rslv.scopes.size() - 1 - i)
			break
		}
	}
}

func (rslv *resolver) beginScope() {
	rslv.scopes.push(make(map[string]bool))
}

func (rslv *resolver) endScope() {
	rslv.scopes.pop()
}

func (rslv *resolver) declare(identifier spec.Token) {
	if rslv.scopes.isEmpty() {
		return
	}
	scope := rslv.scopes.peek()
	if _, contains := scope[identifier.Lexeme]; contains {
		rslv.reportError(identifier, "Already a variable with this name in this scope.")
	}
	scope[identifier.Lexeme] = false
}

func (rslv *resolver) define(identifier spec.Token) {
	if rslv.scopes.isEmpty() {
		return
	}
	scope := rslv.scopes.peek()
	scope[identifier.Lexeme] = true
}

func (rslv *resolver) reportError(token spec.Token, message string) {
	rslv.hadError = true
	rslv.intp.ReportError(token, message)
}

// MARK: - ExprVisitor

func (rslv *resolver) VisitLiteral(le spec.LiteralExpr) (any, error) {
	return nil, nil
}

func (rslv *resolver) VisitGrouping(ge spec.GroupingExpr) (any, error) {
	rslv.resolveExpr(ge.Expr)
	return nil, nil
}

func (rslv *resolver) VisitUnary(ue spec.UnaryExpr) (any, error) {
	rslv.resolveExpr(ue.Expr)
	return nil, nil
}

func (rslv *resolver) VisitBinary(be spec.BinaryExpr) (any, error) {
	rslv.resolveExpr(be.Left)
	rslv.resolveExpr(be.Right)
	return nil, nil
}

func (rslv *resolver) VisitVariable(be spec.VariableExpr) (any, error) {
	if !rslv.scopes.isEmpty() {
		if resolution, contains := rslv.scopes.peek()[be.Identifier.Lexeme]; contains && !resolution {
			rslv.reportError(be.Identifier, "Can't read local variable in its own initializer")
		}
	}
	rslv.resolveLocal(be, be.Identifier)
	return nil, nil
}

func (rslv *resolver) VisitAssignment(ae spec.AssignmentExpr) (any, error) {
	rslv.resolveExpr(ae.Expr);
  rslv.resolveLocal(ae, ae.Identifier);
	return nil, nil
}

func (rslv *resolver) VisitLogical(le spec.LogicalExpr) (any, error) {
	rslv.resolveExpr(le.Left)
	rslv.resolveExpr(le.Right)
	return nil, nil
}

func (rslv *resolver) VisitCall(ce spec.CallExpr) (any, error) {
	rslv.resolveExpr(ce.Callee)
	for _, arg := range ce.Args {
		rslv.resolveExpr(arg)
	}
	return nil, nil
}

func (rslv *resolver) VisitGet(ge spec.GetExpr) (any, error) {
	rslv.resolveExpr(ge.Object)
	return nil, nil
}

func (rslv *resolver) VisitSet(se spec.SetExpr) (any, error) {
	rslv.resolveExpr(se.Object)
	rslv.resolveExpr(se.Value)
	return nil, nil
}

func (rslv *resolver) VisitThis(te spec.ThisExpr) (any, error) {
	if rslv.currentClassType == intp.CtNone {
		rslv.reportError(te.Keyword, "Can't use 'this' outside of a class.")
	}
	rslv.resolveLocal(te, te.Keyword)
	return nil, nil
}

// MARK: - StmtVisitor

func (rslv *resolver) VisitPrint(ps spec.PrintStmt) error {
	rslv.resolveExpr(ps.Expr)
	return nil
}

func (rslv *resolver) VisitExpr(es spec.ExprStmt) error {
	rslv.resolveExpr(es.Expr)
	return nil
}

func (rslv *resolver) VisitDeclare(ds spec.DeclareStmt) error {
	rslv.declare(ds.Identifier)
	if ds.Expr != nil {
		rslv.resolveExpr(ds.Expr)
	}
	rslv.define(ds.Identifier)
	return nil
}

func (rslv *resolver) VisitBlock(bs spec.BlockStmt) error {
	rslv.beginScope()
	rslv.resolveStmts(&bs.Statements)
	rslv.endScope()
	return nil
}

func (rslv *resolver) VisitIf(is spec.IfStmt) error {
	rslv.resolveExpr(is.Condition)
	rslv.resolveStmt(is.Then)
	if is.Else != nil {
		rslv.resolveStmt(is.Else)
	}
	return nil
}

func (rslv *resolver) VisitWhile(ws spec.WhileStmt) error {
	rslv.resolveExpr(ws.Condition)
	rslv.resolveStmt(ws.Body)
	return nil
}

func (rslv *resolver) VisitFunc(fs spec.FuncStmt) error {
	rslv.declare(fs.Name);
  rslv.define(fs.Name);
  rslv.resolveFunction(fs, intp.FtStandalone);
	return nil
}

func (rslv *resolver) VisitReturn(rs spec.ReturnStmt) error {
	if rslv.currentFuncType == intp.FtNone {
		rslv.reportError(rs.Keyword, "Can't return from top-level code")
	}
	if rs.Expr != nil {
		rslv.resolveExpr(rs.Expr)
	}
	return nil
}

func (rslv *resolver) VisitClass(cs spec.ClassStmt) error {
	origClassType := rslv.currentFuncType
	rslv.currentClassType = intp.CtClass
	defer func() { rslv.currentFuncType = origClassType }()

	rslv.declare(cs.Name)
	rslv.define(cs.Name)

	rslv.beginScope()
	rslv.scopes.peek()["this"] = true
	for _, method := range cs.Methods {
		rslv.resolveFunction(method, intp.FtMethod)
	}
	rslv.endScope()
	return nil
}
