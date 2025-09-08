package config

import "time"

type ServerConfig struct {
	Name              string           `mapstructure:"name"`
	User              string           `mapstructure:"user"`
	Host              string           `mapstructure:"host"`
	ConnectionTimeout time.Duration    `mapstructure:"connection_timeout"`
	Auth              ServerAuthConfig `mapstructure:"auth"`
}

type ServerAuthConfig struct {
	Password            string `mapstructure:"password"`
	InteractivePassword bool   `mapstructure:"interactive_password"`
}
