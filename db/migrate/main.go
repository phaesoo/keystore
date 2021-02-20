package main

import (
	"flag"
	"log"

	"github.com/phaesoo/shield/configs"
	"github.com/phaesoo/shield/pkg/db"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	conf := configs.Get().Mysql

	flag.Bool("t", false, "to create test db for integration test")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	test := viper.GetBool("t")

	conn, err := db.NewDB("mysql", conf.DSN())
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	if test {
		database := "shield_test"

		_, err := conn.Exec("DROP DATABASE IF EXISTS ?", database)
		if err != nil {
			log.Fatalf("unable to drop DB `%s`", database)
		}

		_, err = conn.Exec("CREATE DATABASE ?", database)
		if err != nil {
			log.Fatalf("unable to create DB `%s`", database)
		}

		// Patch config with test database
		conf.Database = database

		conn, err = db.NewDB("mysql", conf.DSN())
		defer conn.Close()
		if err != nil {
			panic(err)
		}
	}

}
