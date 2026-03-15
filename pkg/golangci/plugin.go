// pkg/golangci/plugin.go
package golangci

import (
	"github.com/143life/flocci.loglinter/pkg/analyzer"
	"github.com/golangci/plugin-module-register/register" // Официальный регистратор
	"github.com/mitchellh/mapstructure"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("flocciloglint", New)
}

type Plugin struct {
	config *analyzer.Config
}

func New(cfg any) (register.LinterPlugin, error) {
	var settings analyzer.Config
	if err := mapstructure.Decode(cfg, &settings); err != nil {
		return nil, err
	}

	analyzer.SetConfig(&settings)

	return &Plugin{
		config: &settings,
	}, nil
}

// BuildAnalyzers - правильное имя метода для интерфейса register.LinterPlugin
func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

// GetLoadMode возвращает режим загрузки (можно "types" или "syntax")
func (p *Plugin) GetLoadMode() string {
	return "types"
}
