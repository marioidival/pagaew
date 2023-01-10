package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/peterbourgon/ff"
)

func lookupDatabase() string {
	conn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return ""
	}

	return conn
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func migration(dbURL string) error {
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return err
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		return err
	}

	return m.Up()
}
func run() error {
	fs := flag.NewFlagSet("dbmigrate", flag.ExitOnError)

	var (
		databaseURL   = fs.String("database-url", "root:pagaewweagap@/pagaew", "database url")
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		return err
	}

	if err := migration(*databaseURL); err != nil {
		return err
	}

	log.Println("migration cmd done!")

	return nil
}
