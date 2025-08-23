package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type configContextKey struct{}

type Config struct {
	Servers []ServerConfig `mapstructure:"servers"`
	Log     LogConfig      `mapstructure:"log"`
}

func Load() (Config, error) {
	v := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)

	v.SetEnvPrefix("HOMELAB")
	v.SetConfigName("homelab")
	v.AutomaticEnv()

	setConfigPaths(v)
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func setConfigPaths(v *viper.Viper) {
	v.AddConfigPath(".")

	userCfg, err := os.UserConfigDir()
	if err == nil {
		v.AddConfigPath(fmt.Sprintf("%s/homelab", userCfg))
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		v.AddConfigPath(homeDir)
	}
}

func setDefaults(v *viper.Viper) {
	// LogConfig
	v.SetDefault("log.enable", true)
	v.SetDefault("log.level", "info")
}

func ToContext(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configContextKey{}, cfg)
}

func FromContext(ctx context.Context) Config {
	cfg, ok := ctx.Value(configContextKey{}).(Config)
	if !ok {
		panic("config not found in context")
	}

	return cfg
}
