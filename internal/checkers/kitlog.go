package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

type Kitlog struct {
	General
}

func (k Kitlog) FilterExtraArgs(pass *analysis.Pass, args []ast.Expr) []ast.Expr {
	// Check the argument count
	if len(args) > 0 {
		return args[1:]
	}
	return args
}

var _ Checker = (*Zap)(nil)
