package main

import (
	"fmt"

	"github.com/AlexeyGribchenko/selectellint/analyzer"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {

	cfg := analyzer.DefaultConfig()

	if conf != nil {
		if err := mapstructure.Decode(conf, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}
	}

	return []*analysis.Analyzer{
		analyzer.AnalyzerWithConfig(cfg),
	}, nil
}
