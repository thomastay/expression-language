package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
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
//   | KArr
//   | EFieldAccess
//
// Call: Ident ( EList )
//     | E . Ident ( EList )
//
// EFieldAccess : E . Ident
//
// EList : E ',' EList
//       | E
//
// V : Ident | Int | Float | KStr | true | false
//
// KArr  : [ VList ]
// VList : V ',' VList
//       | V

type Expr interface {
	isExpr()
	String() string
}

func (x *EValue) isExpr()       {}
func (x *EBinOp) isExpr()       {}
func (x *EUnOp) isExpr()        {}
func (x *EFieldAccess) isExpr() {}
func (x *EIdxAccess) isExpr()   {}
func (x *ECond) isExpr()        {}
func (x *ECall) isExpr()        {}
func (x *ValueList) isExpr()    {}

type EValue struct {
	Val *lexer.Token
}

func (x *EValue) String() string {
	return x.Val.String()
}

type EBinOp struct {
	Left  Expr
	Op    *lexer.Token
	Right Expr
}

func (x *EBinOp) String() string {
	return fmt.Sprintf("(%s %s %s)", x.Left, string(x.Op.Value), x.Right)
}

type EUnOp struct {
	Op  *lexer.Token
	Val Expr
}

func (x *EUnOp) String() string {
	return fmt.Sprintf("(%s %s)", string(x.Op.Value), x.Val)
}

type EIdxAccess struct {
	Base  Expr
	Index Expr
}

func (x *EIdxAccess) String() string {
	return fmt.Sprintf("%s[%s]", x.Base, x.Index)
}

type EFieldAccess struct {
	Base  Expr
	Field *lexer.Token
}

func (x *EFieldAccess) String() string {
	return fmt.Sprintf("(%s.%s)", x.Base, x.Field)
}

type ECond struct {
	Cond   Expr
	First  Expr
	Second Expr
}

func (x *ECond) String() string {
	return fmt.Sprintf("(%s ? %s : %s)", x.Cond, x.First, x.Second)
}

type ECall struct {
	Base   Expr         // ( @@ "." )?`
	Method *lexer.Token // @Ident`
	Exprs  ExprList     // ( @@ )?`
}

func (x *ECall) String() string {
	if x.Base != nil {
		return fmt.Sprintf("%s.%s(%s)", x.Base, x.Method, x.Exprs)
	}
	return fmt.Sprintf("%s(%s)", x.Method, x.Exprs)
}

type ExprList []Expr

func (es ExprList) String() string {
	if es == nil {
		return ""
	}
	exprs := make([]string, 0, len(es))
	for _, val := range es {
		if val != nil {
			exprs = append(exprs, val.String())
		}
	}
	return strings.Join(exprs, ", ")
}

type ValueList struct {
	Vals []*lexer.Token // "[" ( @@ ( "," @@ )* )? "]"`
}

func (x *ValueList) String() string {
	if x == nil {
		return ""
	}
	exprs := make([]string, len(x.Vals))
	for i, val := range x.Vals {
		exprs[i] = val.String()
	}
	return strings.Join(exprs, ", ")
}
