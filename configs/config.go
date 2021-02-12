package configs

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	fileName  = "config"
	envPrefix = "app"
)

var conf Config

// Config aggregation
type Config struct {
	App   AppConfig
	Mysql MysqlConfig
	Redis RedisConfig
}

// Init is explicit initializer for Config
func init() {
	v := initViper()
	conf = Config{
		App:   appConfig(v),
		Mysql: mysqlConfig(v),
		Redis: redisConfig(v),
	}
}

// Get returns Config object
func Get() Config {
	return conf
}

func initViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName(fileName)

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root := "shield/"
	i := strings.LastIndex(path, "shield/")
	if i != -1 {
		path = path[:i+len(root)]
	}
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// All env vars starts with APP_
	v.AutomaticEnv()
	return v
}
