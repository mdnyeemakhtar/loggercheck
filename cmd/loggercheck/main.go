package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/mdnyeemakhtar/loggercheck"
)

func main() {
	singlechecker.Main(loggercheck.NewAnalyzer())
}
