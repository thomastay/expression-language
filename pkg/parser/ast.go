package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// # Grammmar
// Start : E
// E : V
//   | E op E
//   | uop E
//   | ( E )
//   | E ? E : E
//   | E [ E ]
//   | Call
// 	 | KArr
//
// Call: Ident ( EList )
// 		 | E . Ident ( EList )
// 		 | E . Ident
//
// EList : E ',' EList
// 			 | E
//
// V : Ident | Int | Float | KStr | true | false
//
// KArr  : [ VList ]
// VList : V ',' VList
// 		 	 | V

type Expr interface {
	isExpr()
	String() string
}

func (x *EValue) isExpr()     {}
func (x *EBinOp) isExpr()     {}
func (x *EUnOp) isExpr()      {}
func (x *EIdxAccess) isExpr() {}
func (x *ECond) isExpr()      {}
func (x *Call) isExpr()       {}
func (x *ValueList) isExpr()  {}

type EValue struct {
	val Value
}

func (x *EValue) String() string {
	return x.val.String()
}

type EBinOp struct {
	left  Expr
	op    TokOp
	right Expr
}

func (x *EBinOp) String() string {
	return fmt.Sprintf("(%s %s %s)", x.left, string(x.op), x.right)
}

type EUnOp struct {
	op  TokOp
	val Expr
}

func (x *EUnOp) String() string {
	return fmt.Sprintf("(%s %s)", string(x.op), x.val)
}

type EIdxAccess struct {
	base  Expr
	index Expr
}

func (x *EIdxAccess) String() string {
	return fmt.Sprintf("%s[%s]", x.base, x.index)
}

type ECond struct {
	cond   Expr
	first  Expr
	second Expr
}

func (x *ECond) String() string {
	return fmt.Sprintf("%s ? %s : %s", x.cond, x.first, x.second)
}

type Call struct {
	base   Expr     // ( @@ "." )?`
	method TIdent   // @Ident`
	exprs  ExprList // ( @@ )?`
}

func (x *Call) String() string {
	if x.base != nil {
		return fmt.Sprintf("%s.%s(%s)", x.base, x.method, x.exprs)
	}
	return fmt.Sprintf("%s(%s)", x.method, x.exprs)
}

type ExprList []Expr

func (es ExprList) String() string {
	if es == nil {
		return ""
	}
	exprs := make([]string, len(es))
	for i, val := range es {
		exprs[i] = val.String()
	}
	return strings.Join(exprs, ", ")
}

type ValueList struct {
	vals []Value // "[" ( @@ ( "," @@ )* )? "]"`
}

func (x *ValueList) String() string {
	if x == nil {
		return ""
	}
	exprs := make([]string, len(x.vals))
	for i, val := range x.vals {
		exprs[i] = val.String()
	}
	return strings.Join(exprs, ", ")
}

// Value Union Type
// TODO unify with lexer
type Value interface {
	String() string
}
type TIdent string
type TInt int64
type TFloat float64
type TBool bool
type TKStr string

func (x TIdent) String() string {
	return string(x)
}
func (x TInt) String() string {
	return strconv.FormatInt(int64(x), 10)
}
func (x TFloat) String() string {
	return strconv.FormatFloat(float64(x), 'e', -1, 64)
}
func (x TBool) String() string {
	return string(strconv.FormatBool(bool(x)))
}
func (x TKStr) String() string {
	return string(x)
}
