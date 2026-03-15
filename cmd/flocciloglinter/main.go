package main

import (
	"github.com/143life/flocci.loglinter/pkg/analyzer" // поменяй на свой путь
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	// Устанавливаем конфиг напрямую
	analyzer.SetConfig(&analyzer.Config{
		SensitivePatterns:   []string{"password", "secret"},
		CheckFirstLowercase: true,
		ForbidEmoji:         true,
	})

	// Запускаем анализатор
	singlechecker.Main(analyzer.Analyzer)
}
