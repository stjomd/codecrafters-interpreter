package spec

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"reflect"
)

// MARK: - Expressions

type Expr interface {
	String() string
	Hash() uint64
	Eval(evaluator ExprVisitor[any, error]) (any, error)
}

type ExprVisitor[R any, E error] interface {
	VisitAssignment(assignmentExpr AssignmentExpr) (R, E)
	VisitBinary(binaryExpr BinaryExpr) (R, E)
	VisitCall(callExpr CallExpr) (R, E)
	VisitGrouping(groupingExpr GroupingExpr) (R, E)
	VisitLiteral(literalExpr LiteralExpr) (R, E)
	VisitLogical(logicalExpr LogicalExpr) (R, E)
	VisitUnary(unaryExpr UnaryExpr) (R, E)
	VisitVariable(variableExpr VariableExpr) (R, E)
}

type LiteralExpr struct {
	Value any
}
func (le LiteralExpr) String() string {
	if le.Value == nil {
		return "nil"
	} else if reflect.TypeOf(le.Value).Kind() == reflect.Float64 {
		return float64ToString(le.Value.(float64))
	} else {
		return fmt.Sprint(le.Value)
	}
}
func (le LiteralExpr) Hash() uint64 {
	hash := fnv.New64()
	hash.Write([]byte(le.String()))
	return hash.Sum64()
}
func (le LiteralExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitLiteral(le)
}

type GroupingExpr struct {
	Expr Expr
}
func (ge GroupingExpr) String() string {
	return fmt.Sprintf("(group %v)", ge.Expr)
}
func (ge GroupingExpr) Hash() uint64 {
	hash := fnv.New64()
	hash.Write(bytify(ge.Expr.Hash()))
	return hash.Sum64()
}
func (ge GroupingExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitGrouping(ge)
}

type UnaryExpr struct {
	Opt Token
	Expr Expr
}
func (ue UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", ue.Opt.Lexeme, ue.Expr)
}
func (ue UnaryExpr) Hash() uint64 {
	hash := fnv.New64()
	hash.Write(bytify(ue.Opt.Hash()))
	hash.Write(bytify(ue.Expr.Hash()))
	return hash.Sum64()
}
func (ue UnaryExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitUnary(ue)
}

type BinaryExpr struct {
	Left Expr
	Opt Token
	Right Expr
}
func (be BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", be.Opt.Lexeme, be.Left, be.Right)
}
func (be BinaryExpr) Hash() uint64 {
	hash := fnv.New64()
	hash.Write(bytify(be.Left.Hash()))
	hash.Write(bytify(be.Opt.Hash()))
	hash.Write(bytify(be.Right.Hash()))
	return hash.Sum64()
}
func (be BinaryExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitBinary(be)
}

type VariableExpr struct {
	Identifier Token
	// Random number that identifies a particular occurence of a variable expression (fixes resolvement issues).
	// In the book, the expr->distance map uses the Expr's address as a base for the hash to be used as a key,
	// since the Expr.Variable class in the book does not override hashCode().
	// 
	// In this Go version, both occurences
	// would hash to the same value if this field were not there, which messes up the expr->distance map, as each occurence
	// overwrites the previous one. With this value, which should be chosen at random on each instantiation, an 'address'
	// can be simulated so that the map behaves correctly.
	Occurrence float64
}
func (ve VariableExpr) String() string {
	return fmt.Sprintf("(var %v)", ve.Identifier.Lexeme)
}
func (ve VariableExpr) Hash() uint64 {
	hash := fnv.New64()
	hash.Write(bytify(ve.Identifier.Hash()))
	hash.Write(bytifyFloat64(ve.Occurrence))
	return hash.Sum64()
}
func (ve VariableExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitVariable(ve)
}

type AssignmentExpr struct {
	Identifier Token
	Expr Expr
}
func (ae AssignmentExpr) String() string {
	return fmt.Sprintf("(assign %v %v)", ae.Identifier.Lexeme, ae.Expr)
}
func (ae AssignmentExpr) Hash() uint64 {
	hash := fnv.New64()
	hash.Write(bytify(ae.Identifier.Hash()))
	hash.Write(bytify(ae.Expr.Hash()))
	return hash.Sum64()
}
func (ae AssignmentExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitAssignment(ae)
}

type LogicalExpr struct {
	Left Expr
	Opt Token
	Right Expr
}
func (le LogicalExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", le.Opt.Lexeme, le.Left, le.Right)
}
func (le LogicalExpr) Hash() uint64 {
	hash := fnv.New64()
	hash.Write(bytify(le.Left.Hash()))
	hash.Write(bytify(le.Opt.Hash()))
	hash.Write(bytify(le.Right.Hash()))
	return hash.Sum64()
}
func (le LogicalExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitLogical(le)
}

type CallExpr struct {
	Callee Expr
	Paren Token
	Args []Expr
}
func (ce CallExpr) String() string {
	return fmt.Sprintf("%v(%v)", ce.Callee, ce.Args)
}
func (ce CallExpr) Hash() uint64 {
	hash := fnv.New64()
	hash.Write(bytify(ce.Callee.Hash()))
	hash.Write(bytify(ce.Paren.Hash()))
	for _, arg := range ce.Args {
		hash.Write(bytify(arg.Hash()))
	}
	return hash.Sum64()
}
func (ce CallExpr) Eval(evaluator ExprVisitor[any, error]) (any, error) {
	return evaluator.VisitCall(ce)
}

// MARK: - Helpers

func bytify(hash uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, hash)
	return buf
} 

func bytifyFloat64(num float64) []byte {
	uintRepresentation := math.Float64bits(num);
	return bytify(uintRepresentation)
}
