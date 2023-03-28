package checkers

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/mdnyeemakhtar/loggercheck/internal/checkers/printf"
)

type General struct{}

func (g General) FilterExtraArgs(_ *analysis.Pass, args []ast.Expr) []ast.Expr {
	return args
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
				Message:  fmt.Sprintf("logging message should not use format specifier %q, use fmt.Sprintf to format log message", specifier),
			})

			return // One error diagnostic is enough
		}
	}
}

var _ Checker = (*General)(nil)
