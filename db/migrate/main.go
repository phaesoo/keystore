package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/phaesoo/shield/configs"
	"github.com/phaesoo/shield/db/migrate/migrations"
	"github.com/phaesoo/shield/pkg/db"
	migrate "github.com/phaesoo/sqlx-migrate"
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

	// Patch config if test
	if test {
		database := "shield_test"
		conf.Database = database

		_, err := conn.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", database))
		if err != nil {
			log.Fatalf("unable to drop DB `%s`: %s", database, err.Error())
		}
		_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", database))
		if err != nil {
			log.Fatalf("unable to create DB `%s`", database)
		}

		conn, err = db.NewDB("mysql", conf.DSN())
		defer conn.Close()
		if err != nil {
			panic(err)
		}
	}

	m := migrate.New(conn, []migrate.Migration{
		migrations.InitTables,
	})
	if err := m.Migrate(); err != nil {
		log.Printf("Migration failed: %s", err.Error())
		if err := m.Rollback(); err != nil {
			log.Printf("Failed to rollback last migration: %s", err.Error())
		}
		log.Printf("Success to rollback last migration")
	}
}
