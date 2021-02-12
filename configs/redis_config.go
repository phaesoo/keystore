package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host     string
	Port     string
	Database int
	Logging  bool
}

func redisConfig(v *viper.Viper) RedisConfig {
	return RedisConfig{
		Host:     v.GetString("redis.host"),
		Port:     v.GetString("redis.port"),
		Database: v.GetInt("redis.database"),
		Logging:  v.GetBool("redis.logging"),
	}
}

func (c RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
