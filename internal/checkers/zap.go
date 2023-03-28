package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Zap struct {
	General
}

func (z Zap) FilterExtraArgs(pass *analysis.Pass, args []ast.Expr) []ast.Expr {
	// Check the argument count
	filtered := make([]ast.Expr, 0, len(args))
	for _, arg := range args {
		// Skip any zapcore.Field we found
		switch arg := arg.(type) {
		case *ast.CallExpr, *ast.Ident:
			typ := pass.TypesInfo.TypeOf(arg)
			switch typ := typ.(type) {
			case *types.Named:
				obj := typ.Obj()
				// This is a strongly-typed field. Consume it and move on.
				// Actually it's go.uber.org/zap/zapcore.Field, however for simplicity
				// we don't check the import path
				if obj != nil && obj.Name() == "Field" {
					continue
				}
			default:
				// pass
			}
		}

		filtered = append(filtered, arg)
	}

	return filtered
}

var _ Checker = (*Zap)(nil)
