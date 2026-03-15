// pkg/golangci/plugin.go
package golangci

import (
	"fmt"
	"os"

	"github.com/143life/flocci.loglinter/pkg/analyzer"
	"github.com/golangci/plugin-module-register/register" // Официальный регистратор
	"github.com/mitchellh/mapstructure"
	"golang.org/x/tools/go/analysis"
)

func init() {
	// Регистрируем плагин под именем "flocciloglint"
	fmt.Fprintf(os.Stderr, ">>> PLUGIN INIT CALLED\n")
	register.Plugin("flocciloglint", New)
}

type Plugin struct {
	config *analyzer.Config
}

func New(cfg any) (register.LinterPlugin, error) {
	fmt.Fprintf(os.Stderr, "=== PLUGIN: New called with config: %+v ===\n", cfg)

	var settings analyzer.Config
	if err := mapstructure.Decode(cfg, &settings); err != nil {
		fmt.Fprintf(os.Stderr, "=== PLUGIN: mapstructure error: %v ===\n", err)
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "=== PLUGIN: decoded settings: %+v ===\n", settings)
	analyzer.SetConfig(&settings)

	return &Plugin{
		config: &settings,
	}, nil
}

// BuildAnalyzers - правильное имя метода для интерфейса register.LinterPlugin
func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	fmt.Fprintf(os.Stderr, ">>> BuildAnalyzers called\n")
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

// GetLoadMode возвращает режим загрузки (можно "types" или "syntax")
func (p *Plugin) GetLoadMode() string {
	return "types"
}
