package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marioidival/pagaew/cmd/server/handlers"
	"github.com/peterbourgon/ff"
)

func main() {

	// TODO: criar um pkg do server
	// TODO: criar um pkg para db connection
	// TODO: pegar as infos via env (dotenv)

	if err := run(); err != nil {
		log.Fatal(err.Error())
	}

}

func run() error {
	fs := flag.NewFlagSet("server", flag.ExitOnError)

	var (
		databaseURL   = fs.String("database-url", "", "database url")
		webserverAddr = fs.String("addr", ":3000", "webserver addr - default :3000")
	)

	fmt.Println(databaseURL)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		return err
	}

	// server setup
	mux := handlers.Setup()

	server := &http.Server{
		Addr: *webserverAddr,
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
