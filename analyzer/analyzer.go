package analyzer

import "golang.org/x/tools/go/analysis"

func AnalyzerWithConfig(cfg Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "selectellint",
		Doc:  "Linter for logging messages to ensure they follow best practices: lowercase start, no invalid symbols, no sensitive data",
		Run:  RunWithConfig(cfg),
	}
}
