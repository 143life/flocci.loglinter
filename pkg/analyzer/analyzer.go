package analyzer

import (
	"flag"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:  "flocciloglint",
	Doc:   "Checks log messages for sensitive data and formatting issues",
	Run:   run,
	Flags: flags(),
}

func flags() flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.StringVar(&configFlag, "config", "", "path to configuration file (YAML or JSON)")
	return *fs
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Determine function name without using types.
			var funcName string
			switch fun := call.Fun.(type) {
			case *ast.Ident:
				funcName = fun.Name
			case *ast.SelectorExpr:
				if pkg, ok := fun.X.(*ast.Ident); ok {
					funcName = pkg.Name + "." + fun.Sel.Name
				} else {
					funcName = fun.Sel.Name
				}
			default:
				return true
			}

			// Check if it's a function from the "log" package.
			if strings.HasPrefix(funcName, "log.") {
				if len(call.Args) == 0 {
					return true
				}
				firstArg := call.Args[0]
				lit, ok := firstArg.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}

				msg := strings.Trim(lit.Value, "\"`")
				cfg, err := getConfig()
				if err != nil {
					cfg = DefaultConfig()
				}

				// Check for sensitive patterns.
				if len(cfg.SensitivePatterns) > 0 {
					lowerMsg := strings.ToLower(msg)
					for _, pattern := range cfg.SensitivePatterns {
						if strings.Contains(lowerMsg, strings.ToLower(pattern)) {
							pass.Reportf(firstArg.Pos(), "potential sensitive data: %q", pattern)
							break
						}
					}
				}

				// Check first character case.
				if cfg.CheckFirstLowercase && len(msg) > 0 && !isLowercase(msg) {
					pass.Reportf(firstArg.Pos(), "log message should start with lowercase")
				}

				// Check for emoji.
				if cfg.ForbidEmoji && containsEmoji(msg) {
					pass.Reportf(firstArg.Pos(), "log message contains emoji")
				}

				// Check for non-ASCII characters (only English).
				if cfg.AllowOnlyASCII && !isASCII(msg) {
					pass.Reportf(firstArg.Pos(), "log message contains non-ASCII characters")
				}
			}
			return true
		})
	}
	return nil, nil
}
