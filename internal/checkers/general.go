package checkers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/timonwong/loggercheck/internal/checkers/printf"
)

type General struct{}

func (g General) FilterKeyAndValues(_ *analysis.Pass, keyAndValues []ast.Expr) []ast.Expr {
	return keyAndValues
}

func (g General) CheckPrintfLikeSpecifier(pass *analysis.Pass, args []ast.Expr) {
	for _, arg := range args {
		format, ok := extractValueFromStringArg(pass, arg)
		if !ok {
			continue
		}

		if specifier, ok := printf.IsPrintfLike(format); ok {
			pass.Report(analysis.Diagnostic{
				Pos:      arg.Pos(),
				End:      arg.End(),
				Category: DiagnosticCategory,
				Message:  fmt.Sprintf("logging message should not use format specifier %q", specifier),
			})

			return // One error diagnostic is enough
		}
	}
}

var _ Checker = (*General)(nil)
