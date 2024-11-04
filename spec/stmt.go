package spec

type Stmt interface {
	Exec(executor StmtVisitor[error]) error
}

type StmtVisitor[R any] interface {
	VisitPrint(printStmt PrintStmt) R
	VisitExpr(exprStmt ExprStmt) R
	VisitDeclare(declareStmt DeclareStmt) R
	VisitBlock(blockStmt BlockStmt) R
	VisitIf(ifStmt IfStmt) R
}

type PrintStmt struct {
	Expr Expr
}
func (ps PrintStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitPrint(ps)
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
