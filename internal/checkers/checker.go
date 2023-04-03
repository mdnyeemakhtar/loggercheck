package checkers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Config struct {
	RequireStringKey bool
	NoPrintfLike     bool
}

type CallContext struct {
	Expr      *ast.CallExpr
	Func      *types.Func
	Signature *types.Signature
}

type Checker interface {
	FilterExtraArgs(pass *analysis.Pass, keyAndValues []ast.Expr) []ast.Expr
	CheckPrintfLikeSpecifier(pass *analysis.Pass, args []ast.Expr)
}

func ExecuteChecker(c Checker, pass *analysis.Pass, call CallContext, cfg Config) {
	c.CheckPrintfLikeSpecifier(pass, call.Expr.Args)
	params := call.Signature.Params()
	nparams := params.Len() // variadic => nonzero
	startIndex := nparams - 1

	lastArg := params.At(nparams - 1)
	iface, ok := lastArg.Type().(*types.Slice).Elem().(*types.Interface)
	if !ok || !iface.Empty() {
		return // final (args) param is not ...interface{}
	}

	extraArgs := c.FilterExtraArgs(pass, call.Expr.Args[startIndex:])

	if len(extraArgs) > 0 {
		firstArg := extraArgs[0]
		lastArg := extraArgs[len(extraArgs)-1]
		pass.Report(analysis.Diagnostic{
			Pos:      firstArg.Pos(),
			End:      lastArg.End(),
			Category: DiagnosticCategory,
			Message:  "additional arguments passed for logging, use fmt.Sprintf to format log message",
		})
	}
}
