package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/peterbourgon/ff"

	"github.com/marioidival/pagaew/cmd/server/handlers"
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
		webserverAddr = fs.String("addr", ":3000", "webserver addr - default :3000")
		environment = fs.String("ENVIRONMENT", "test", "environment of system")
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

	// server setup
	mux := handlers.Setup(dbc, *environment == "prod")

	server := &http.Server{
		Addr:    *webserverAddr,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("startup server addr", *webserverAddr)
		serverErrors <- server.ListenAndServe()
	}()

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case serverError := <-serverErrors:
		return serverError

	case sig := <-quit:
		log.Println("closing the server ->", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if shutdownError := server.Shutdown(ctx); shutdownError != nil {
			defer func() {
				if closeErr := server.Close(); closeErr != nil {
					log.Fatalln("Could not close the server", closeErr)
				}
			}()
			log.Fatalln("Could not gracefully shutdown the server")
		}
		close(done)
	case <-done:
		return nil
	}

	return nil
}
