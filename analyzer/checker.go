package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func RunWithConfig(cfg Config) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {

		c := Checker{cfg: &cfg}
		for _, file := range pass.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				if len(call.Args) == 0 {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if c.isLoggerFunc(sel.Sel.Name) {
					for _, arg := range call.Args {

						switch expr := arg.(type) {
						case *ast.BasicLit:
							if cfg.CheckCapital {
								c.checkCapital(pass, expr)
							}
							if cfg.CheckInvalid {
								c.checkInvalid(pass, expr)
							}
						default:
							c.basePos = arg.Pos()
							c.checkExpression(pass, arg)
						}
					}
				}
				return true
			})
		}
		return nil, nil
	}

}

type Checker struct {
	basePos token.Pos
	cfg     *Config
}

func (c *Checker) checkCapital(pass *analysis.Pass, expr *ast.BasicLit) {
	if expr.Kind != token.STRING {
		return
	}

	clearText := strings.Trim(expr.Value, "\"")
	if len(clearText) > 0 && c.hasFirstCapital(clearText) {
		newText := c.fixCapital(clearText)
		pass.Report(analysis.Diagnostic{
			Pos:     expr.ValuePos + 1,
			End:     expr.End(),
			Message: "message should start with lowercase",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "convert first letter to lowercase",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     expr.ValuePos + 1,
							End:     expr.End(),
							NewText: []byte(newText),
						},
					},
				},
			},
		})
	}
}

func (c *Checker) checkInvalid(pass *analysis.Pass, expr *ast.BasicLit) {
	if expr.Kind != token.STRING {
		return
	}

	if ok, pos := c.hasInvalidSymbol(expr.Value); ok {
		newText := c.fixInvalid(expr.Value)
		pass.Report(analysis.Diagnostic{
			Pos:     expr.ValuePos + token.Pos(pos) + 1,
			End:     expr.End(),
			Message: "invalid symbol",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					TextEdits: []analysis.TextEdit{
						{
							Pos:     expr.ValuePos + 1,
							End:     expr.End() - 1,
							NewText: []byte(newText),
						},
					},
				},
			},
		})
	}
}

func (c Checker) checkSensetive(pass *analysis.Pass, expr *ast.BasicLit) {
	if ok, pos := c.hasSensetiveData(expr.Value); ok {
		pass.Reportf(expr.ValuePos+token.Pos(pos), "sensetive data")
	}
}

func (c Checker) checkVariable(pass *analysis.Pass, expr *ast.Ident) {
	if ok, _ := c.hasSensetiveData(expr.Name); ok {
		pass.Reportf(expr.NamePos, "sensetive data {%s}", expr.Name)
	}
}

func (c Checker) checkExpression(pass *analysis.Pass, arg ast.Expr) {
	switch expr := arg.(type) {
	case *ast.StarExpr:
		c.checkExpression(pass, expr.X)
	case *ast.IndexExpr:
		c.checkExpression(pass, expr.X)
	case *ast.ParenExpr:
		c.checkExpression(pass, expr.X)
	case *ast.SliceExpr:
		c.checkExpression(pass, expr.X)
	case *ast.BinaryExpr:
		c.checkExpression(pass, expr.X)
		c.checkExpression(pass, expr.Y)
	case *ast.SelectorExpr:
		c.checkExpression(pass, expr.Sel)
	case *ast.CallExpr:
		c.checkExpression(pass, expr.Fun)
		for _, arg := range expr.Args {
			c.checkExpression(pass, arg)
		}
	case *ast.Ident:
		c.checkVariable(pass, expr)
	case *ast.BasicLit:
		if expr.Pos()-c.basePos == 0 {
			if c.cfg.CheckCapital {
				c.checkCapital(pass, expr)
			}
		}
		if c.cfg.CheckInvalid {
			c.checkInvalid(pass, expr)
		}
		if c.cfg.CheckSensitive {
			c.checkSensetive(pass, expr)
		}
	}
}

func (c Checker) isLoggerFunc(name string) bool {
	loggerFuncNames := []string{
		"warn", "info", "panic", "error", "log", "debug",
	}
	name = strings.ToLower(name)
	flag := false
	for _, funcName := range loggerFuncNames {
		if strings.Contains(name, funcName) {
			flag = true
			break
		}
	}
	return flag
}

func (c Checker) hasFirstCapital(text string) bool {
	if len(text) == 0 {
		return false
	}
	return 'A' <= text[0] && text[0] <= 'Z'
}

func (c Checker) hasInvalidSymbol(text string) (bool, int) {
	clearText := strings.Trim(text, "\"")
	for i, r := range clearText {
		if !c.isSymbolValid(r) {
			return true, i
		}
	}
	return false, 0
}

func (c Checker) hasSensetiveData(text string) (bool, int) {
	sensetive := []string{"password", "token", "login", "email", "id", "api", "credential"}
	lowerText := strings.ToLower(text)
	for _, word := range sensetive {
		if strings.Contains(lowerText, word) {
			return true, strings.Index(lowerText, word)
		}
	}
	return false, 0
}

func (c Checker) isSymbolValid(r rune) bool {
	return unicode.In(r, unicode.Latin) || unicode.IsDigit(r) || r == rune(' ') || strings.ContainsRune(c.cfg.AllowedSymbols, r)
}

func (c Checker) fixCapital(text string) string {
	if len(text) == 0 {
		return ""
	}
	return strings.ToLower(text[:1]) + text[1:]
}

func (c Checker) fixInvalid(text string) string {
	var newText strings.Builder
	newText.Grow(len(text))

	for i, r := range text {
		if strings.ContainsRune(c.cfg.ReplaceWithSpace, r) {
			if i != 0 && i != len(text)-1 {
				newText.WriteRune(' ')
			}
			continue
		}
		if !c.isSymbolValid(r) {
			continue
		}
		newText.WriteRune(r)
	}
	return newText.String()
}
