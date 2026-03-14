package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "flocciloglinter",
	Doc:      "blahblahblah",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		typ := pass.TypesInfo.TypeOf(call.Fun)

		if typ == nil {
			return
		}

		_, ok := typ.(*types.Signature)
		if !ok {
			return
		}

		funObj, ok := call.Fun.(*ast.Ident)
		if !ok {
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}

			obj := pass.TypesInfo.ObjectOf(sel.Sel)
			if obj == nil {
				return
			}

			if fn, ok := obj.(*types.Func); ok {
				checkLogCall(pass, call, fn)
			}
		} else {
			obj := pass.TypesInfo.ObjectOf(funObj)
			if fn, ok := obj.(*types.Func); ok {
				checkLogCall(pass, call, fn)
			}
		}
	})

	return nil, nil
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr, fn *types.Func) {
	pkgName := fn.Pkg().Name()
	funcName := fn.Name()

	if pkgName == "log" && (funcName == "Print" || funcName == "Printf" || funcName == "Println" ||
		funcName == "Fatal" || funcName == "Fatalf" || funcName == "Fatalln" ||
		funcName == "Panic" || funcName == "Panicf" || funcName == "Panicln") {
		extractAndCheckMessage(pass, call)
		return
	}

	if pkgName == "slog" && (funcName == "Info" || funcName == "Debug" || funcName == "Warn" || funcName == "Error") {
		extractAndCheckMessage(pass, call)
		return
	}
}

func extractAndCheckMessage(pass *analysis.Pass, call *ast.CallExpr) {
	pass.Reportf(call.Pos(), "found a function call")
	if len(call.Args) == 0 {
		return
	}
	firstArg := call.Args[0]
	lit, ok := firstArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}
	msg := strings.Trim(lit.Value, "\"")

	if !isLowercase(msg) {
		pass.Reportf(firstArg.Pos(), "log message should start with lowercase")
	}
	// TODO: остальные требования
}

func isLowercase(msg string) bool {
	if len(msg) == 0 {
		return true
	}
	r, _ := utf8.DecodeRuneInString(msg)
	return unicode.IsLower(r)
}
