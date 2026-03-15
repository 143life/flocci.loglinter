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

			var isLog bool
			msgArgIndex := 0

			if strings.HasPrefix(funcName, "log.") {
				isLog = true
			} else if strings.HasPrefix(funcName, "slog.") {
				base := strings.TrimPrefix(funcName, "slog.")
				switch base {
				case "Info", "Debug", "Warn", "Error":
					isLog = true
				case "InfoContext", "DebugContext", "WarnContext", "ErrorContext":
					isLog = true
					msgArgIndex = 1
				}
			}
			// TODO: добавить поддержку zap

			if !isLog {
				return true
			}

			if len(call.Args) <= msgArgIndex {
				return true
			}
			msgArg := call.Args[msgArgIndex]
			lit, ok := msgArg.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}

			msg := strings.Trim(lit.Value, "\"`")
			cfg, err := getConfig()
			if err != nil {
				cfg = DefaultConfig()
			}

			if len(cfg.SensitivePatterns) > 0 {
				lowerMsg := strings.ToLower(msg)
				for _, pattern := range cfg.SensitivePatterns {
					if strings.Contains(lowerMsg, strings.ToLower(pattern)) {
						pass.Reportf(msgArg.Pos(), "potential sensitive data: %q", pattern)
						break
					}
				}
			}

			if cfg.CheckFirstLowercase && len(msg) > 0 && !isLowercase(msg) {
				pass.Reportf(msgArg.Pos(), "log message should start with lowercase")
			}

			if cfg.ForbidEmoji && containsEmoji(msg) {
				pass.Reportf(msgArg.Pos(), "log message contains emoji")
			}

			if cfg.ForbidSpecialChars && containsSpecialChars(msg) {
				pass.Reportf(msgArg.Pos(), "log message contains special characters")
			}

			if cfg.AllowOnlyASCII && !isASCII(msg) {
				pass.Reportf(msgArg.Pos(), "log message contains non-ASCII characters")
			}

			return true
		})
	}
	return nil, nil
}
