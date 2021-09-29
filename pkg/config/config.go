/**
 * @File: config.go
 * @Author: hsien
 * @Description:
 * @Date: 9/17/21 6:16 PM
 */

package config

import (
	"github.com/BurntSushi/toml"
)

var (
	Debug bool
)

func DefaultConfig() *Config {
	return &Config{
		Server: &Server{
			Addr:  ":8088",
			Debug: false,
		},
		Log: &Log{
			LogLevel: "DEBUG",
			LogFile:  "",
		},
	}
}

func LoadFromFile(path string) (*Config, error) {
	cfg := new(Config)
	if path != "" {
		_, err := toml.DecodeFile(path, cfg)
		if err != nil {
			return cfg, err
		}
	}
	Debug = cfg.Server.Debug
	return cfg, nil
}
