package checkers

import (
	"go/ast"
	"go/constant"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const (
	DiagnosticCategory = "logging"
)

// extractValueFromStringArg returns true if the argument is a string type (literal or constant).
func extractValueFromStringArg(pass *analysis.Pass, arg ast.Expr) (value string, ok bool) {
	if typeAndValue, ok := pass.TypesInfo.Types[arg]; ok {
		if typ, ok := typeAndValue.Type.(*types.Basic); ok && typ.Kind() == types.String && typeAndValue.Value != nil {
			return constant.StringVal(typeAndValue.Value), true
		}
	}

	return "", false
}
