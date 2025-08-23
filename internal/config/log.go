package config

type LogConfig struct {
	Enable bool   `mapstructure:"enable"`
	Level  string `mapstructure:"level"`
}
