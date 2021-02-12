package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type MysqlConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	Logging  bool
}

func mysqlConfig(v *viper.Viper) MysqlConfig {
	return MysqlConfig{
		Host:     v.GetString("mysql.host"),
		Port:     v.GetInt("mysql.port"),
		Database: v.GetString("mysql.database"),
		User:     v.GetString("mysql.user"),
		Password: v.GetString("mysql.password"),
		Logging:  v.GetBool("mysql.logging"),
	}
}

func (c MysqlConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database)
}
