package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/peterbourgon/ff"

	"github.com/marioidival/pagaew/pkg/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	fs := flag.NewFlagSet("server", flag.ExitOnError)

	var (
		databaseURL   = fs.String("database-url", "root:pagaewweagap@/pagaew", "database url")
		interval = fs.Int("INTERVAL", 10, "interval in seconds")
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		return err
	}

	dbc, err := database.Open(context.Background(), *databaseURL)
	if err != nil {
		return err
	}
	if err := dbc.Ping(); err != nil {
		return err
	}
	defer dbc.Close()

	for {
		rows, err := dbc.Query(context.Background(), "SELECT debt_id FROM log_invoice WHERE status = ? FOR UPDATE SKIP LOCKED", "PENDING")
		if err != nil {
			log.Println("log invoice failed to get PENDING to proceed", err)
			break
		}

		for rows.Next() {
			var debtId string
			err := rows.Scan(&debtId)
			if err != nil {
				log.Println("failed to scan values")
				continue
			}

			// simulate a job to sent an email
			_, err = dbc.Exec(context.Background(), "UPDATE log_invoice SET status = ? WHERE debt_id = ?", "EMAIL_SENT", debtId)
			if err != nil {
				return err
			}
		}
		log.Println("N jobs was updated")

		time.Sleep(time.Duration(*interval * int(time.Second)))
	}

	return nil
}
