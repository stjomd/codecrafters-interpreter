package spec

import "fmt"

type Stmt interface {
	Exec(executor StmtVisitor[error]) error
}

type StmtVisitor[R any] interface {
	VisitBlock(blockStmt BlockStmt) R
	VisitDeclare(declareStmt DeclareStmt) R
	VisitExpr(exprStmt ExprStmt) R
	VisitFunc(funcStmt FuncStmt) R
	VisitIf(ifStmt IfStmt) R
	VisitPrint(printStmt PrintStmt) R
	VisitWhile(whileStmt WhileStmt) R
}

type PrintStmt struct {
	Expr Expr
}
func (ps PrintStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitPrint(ps)
}
func (ps PrintStmt) String() string {
	return fmt.Sprintf("(print %v)", ps.Expr)
}

type ExprStmt struct {
	Expr Expr
}
func (es ExprStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitExpr(es)
}

type DeclareStmt struct {
	Identifier Token
	Expr Expr
}
func (ds DeclareStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitDeclare(ds)
}

type BlockStmt struct {
	Statements []Stmt
}
func (bs BlockStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitBlock(bs)
}

type IfStmt struct {
	Condition Expr
	Then Stmt
	Else Stmt
}
func (is IfStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitIf(is)
}

type WhileStmt struct {
	Condition Expr
	Body Stmt
}
func (ws WhileStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitWhile(ws)
}

type FuncStmt struct {
	Name Token
	Params []Token
	Body BlockStmt
}
func (fs FuncStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitFunc(fs)
}
