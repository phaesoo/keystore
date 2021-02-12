package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// AppConfig is config struct for app
type AppConfig struct {
	Host    string
	Port    int
	Profile bool
	Metrics bool
}

func appConfig(v *viper.Viper) AppConfig {
	return AppConfig{
		Host:    v.GetString("app.host"),
		Port:    v.GetInt("app.port"),
		Profile: v.GetBool("app.profile"),
		Metrics: v.GetBool("app.metrics"),
	}
}

func (c *AppConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
