// pkg/golangci/config.go
package golangci

type Settings struct {
	SensitivePatterns []string `mapstructure:"sensitive"`
	NoEmoji           bool     `mapstructure:"no-emoji"`
	FirstLowercase    bool     `mapstructure:"first-lowercase"`
	//MaxLength         int      `mapstructure:"max-length"`
}

// В функции New(cfg any) ты приводишь cfg к этому типу:
// var s Settings
// if err := mapstructure.Decode(cfg, &s); err != nil { ... }
