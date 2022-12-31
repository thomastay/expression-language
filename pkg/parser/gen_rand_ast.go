package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

type Rand struct {
	x uint32
}

func (r *Rand) randU32() uint32 {
	// from https://github.com/valyala/fastrand/blob/master/fastrand.go
	// See https://en.wikipedia.org/wiki/Xorshift
	x := r.x
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.x = x
	return x
}

// Mostly for testing
func GenRandomAST(seed uint32, maxDepth uint) Expr {
	rand := Rand{seed + 1}
	return genRandomASTRec(&rand, maxDepth)
}

func genRandomASTRec(rand *Rand, depthLeft uint) Expr {
	seed := rand.randU32()
	nAstNodesGenTypes := uint32(numASTNodeTypes)
	if depthLeft == 0 {
		// reset to EValue
		seed = seed / nAstNodesGenTypes * nAstNodesGenTypes
	}
	maxArraySize := 5
	value := seed
	nodeType := seed % nAstNodesGenTypes
	switch nodeType {
	case 0:
		// EValue
		value /= nAstNodesGenTypes
		switch value % 6 { // ident, true, false, int, float, string
		case 0:
			return &EValue{
				Val: &lexer.Token{
					Type:  TokIdent,
					Value: seedToIdent(seed),
				},
			}
		case 1:
			return &EValue{
				Val: &lexer.Token{
					Type:  TokBool,
					Value: "true",
				},
			}
		case 2:
			return &EValue{
				Val: &lexer.Token{
					Type:  TokBool,
					Value: "false",
				},
			}
		case 3:
			return &EValue{
				Val: &lexer.Token{
					Type:  TokInt,
					Value: fmt.Sprint(seed),
				},
			}
		case 4:
			return &EValue{
				Val: &lexer.Token{
					Type:  TokFloat,
					Value: fmt.Sprint(seed),
				},
			}
		case 5:
			return &EValue{
				Val: &lexer.Token{
					Type:  TokSingleString,
					Value: fmt.Sprint(seed),
				},
			}
		}
	case 1:
		// EUnOp
		value /= nAstNodesGenTypes
		opStr := unaryOps[value%uint32(len(unaryOps))]
		return &EUnOp{
			Op: &lexer.Token{
				Type:  TokOp,
				Value: opStr,
			},
			Val: genRandomASTRec(rand, depthLeft-1),
		}
	case 2:
		// EBinOp
		value /= nAstNodesGenTypes
		opStr := binaryOps[value%uint32(len(binaryOps))]

		return &EBinOp{
			Left:  genRandomASTRec(rand, depthLeft-1),
			Right: genRandomASTRec(rand, depthLeft-1),
			Op: &lexer.Token{
				Type:  TokOp,
				Value: opStr,
			},
		}
	case 3:
		// EFieldAccess
		return &EFieldAccess{
			Base: genRandomASTRec(rand, depthLeft-1),
			Field: &lexer.Token{
				Type:  TokIdent,
				Value: seedToIdent(seed),
			},
		}
	case 4:
		return &EIdxAccess{
			Base:  genRandomASTRec(rand, depthLeft-1),
			Index: genRandomASTRec(rand, depthLeft-1),
		}
	case 5:
		return &ECond{
			Cond:   genRandomASTRec(rand, depthLeft-1),
			First:  genRandomASTRec(rand, depthLeft-1),
			Second: genRandomASTRec(rand, depthLeft-1),
		}
	case 6:
		return &ECall{
			Base: genRandomASTRec(rand, depthLeft-1),
			Method: &lexer.Token{
				Type:  TokIdent,
				Value: seedToIdent(seed),
			},
			Exprs: []Expr{
				genRandomASTRec(rand, depthLeft-1),
			},
		}
	case 7:
		arrLen := rand.randU32() % uint32(maxArraySize)
		arr := make([]Expr, arrLen)
		for i := 0; i < int(arrLen); i++ {
			arr[i] = genRandomASTRec(rand, depthLeft-1)
		}
		return (*EArray)(&arr)
	}
	panic("Not impl")
}

var unaryOps = []string{
	"+",
	"-",
	"not",
}

var binaryOps = []string{
	"+",
	"-",
	"*",
	"/",
	"//",
	"%",
	"**",
	"<",
	">",
	">=",
	"<=",
	"==",
	"!=",
	"and",
	"or",
}

var identsSeeded = []string{
	"a",
	"b",
	"c",
	"d",
	"e",
	"f",
	"foo",
	"s",
	"fizz",
	"buzz",
	"fizzbuzz",
	"emptyObj",
	"fooObj",
	"bar",
	"baz",
	"foobar",
	"z",
	"ba",
	"vv",
	"null",
}

func seedToIdent(seed uint32) string {
	// No point in creating an ident that doesn't exist, doesn't give useful info since VM stops immediately
	return identsSeeded[seed%uint32(len(identsSeeded))]
}
