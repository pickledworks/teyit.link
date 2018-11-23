// teyit.link
package main

import (
	"context"
	"github.com/noddigital/teyit.link/database"
	"github.com/noddigital/teyit.link/handlers"
	"github.com/noddigital/teyit.link/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//go:generate statik -src=./public/

func main() {
	config := utils.InitConfig()

	db := database.InitDB(config.DbDialect, config.DbUri)
	database.Migrate(db)
	defer db.Close()

	srv := &http.Server{
		Addr: config.ServerAddr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		// Pass our instance of gorilla/mux in.
		// Handler:      CSRF(handlers.CreateRoutes()),
		Handler: handlers.CreateRoutes(),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Panic(err)
		} else {
			log.Printf("Server started")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), config.GracefulShutdown)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
