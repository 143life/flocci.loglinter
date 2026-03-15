package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SensitivePatterns   []string `mapstructure:"sensitive_patterns"`
	UseRegex            bool     `mapstructure:"use_regex"`
	CheckFirstLowercase bool     `mapstructure:"check_first_lowercase"`
	ForbidEmoji         bool     `mapstructure:"forbid_emoji"`
	//MaxLength           int      `mapstructure:"max_length"`
	AllowOnlyASCII bool `mapstructure:"allow_only_ascii"`
}

var (
	mu             sync.RWMutex
	configFlag     string // path to file config
	configInstance *Config
)

func DefaultConfig() *Config {
	return &Config{
		SensitivePatterns:   []string{"password", "token", "secret", "key"},
		CheckFirstLowercase: true,
		ForbidEmoji:         true,
		//MaxLength:           0,
		AllowOnlyASCII: false,
	}
}

func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file (as YAML or JSON): %w", err)
		}
	}
	return &cfg, nil
}

// SetConfig устанавливает конфигурацию для анализатора
func SetConfig(cfg *Config) {
	mu.Lock()
	defer mu.Unlock()
	configInstance = cfg
}

// getConfig возвращает текущую конфигурацию (с дефолтом)
func getConfig() (*Config, error) {
	mu.RLock()
	defer mu.RUnlock()

	if configInstance != nil {
		return configInstance, nil
	}

	if configFlag != "" {
		cfg, err := LoadConfigFromFile(configFlag)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
		return cfg, nil
	}

	// Если конфигурация не установлена, возвращаем дефолтную
	return DefaultConfig(), nil
}
