package spec

type Stmt interface {
	Exec(executor StmtVisitor[error]) error
}

type StmtVisitor[R any] interface {
	VisitPrint(printStmt PrintStmt) R
}

type PrintStmt struct {
	Expr Expr
}
func (ps PrintStmt) Exec(executor StmtVisitor[error]) error {
	return executor.VisitPrint(ps)
}
