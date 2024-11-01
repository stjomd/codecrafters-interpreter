package spec

type Stmt interface {
	Exec(executor StmtVisitor[error]) error
}

type StmtVisitor[R any] interface {
	VisitPrint(printStmt PrintStmt) R
	VisitExpr(exprStmt ExprStmt) R
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
