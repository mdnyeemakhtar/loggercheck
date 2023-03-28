package checkers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

type Kitlog struct {
	General
}

func (k Kitlog) FilterKeyAndValues(pass *analysis.Pass, keyAndValues []ast.Expr) []ast.Expr {
	// Check the argument count
	if len(keyAndValues) > 0 {
		return keyAndValues[1:]
	}
	return keyAndValues
}

var _ Checker = (*Zap)(nil)
