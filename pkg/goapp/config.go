package goapp

import (
	"strings"

	"github.com/spf13/viper"
)

// Config is a viper based application config
var Config = viper.New()

// Sub extracts Sub config from viper using env variables
func Sub(config *viper.Viper, name string) *viper.Viper {
	res := config.Sub(name)
	return res
}

// InitEnv initializes viper for environment variables
func InitEnv(config *viper.Viper) {
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
}
