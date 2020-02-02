package config

import (
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/env"
	// "github.com/micro/go-plugins/config/source/configmap"
)

// NewConfig returns config with env and k8s configmap setup
func NewConfig(opts ...config.Option) config.Config {
	cfg, _ := config.NewConfig()
	cfg.Load(
		env.NewSource(),
	)
	return cfg
}

// Config global config
var Config config.Config

func init() {
	Config = NewConfig()
	// fmt.Println(Config)
	// Config = NewConfig()
}
