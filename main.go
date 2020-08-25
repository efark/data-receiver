/*
The execution of the application starts here, in the 'main' package.

*/
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/efark/data-receiver/webserver"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.Info("Starting webserver.")

	var cfgFilepath, cfgInline string
	flag.StringVar(&cfgFilepath, "config", "", "Configuration file path (json or yaml).")
	flag.StringVar(&cfgInline, "inline-config", "", "Inline Json Configuration.")
	flag.Parse()

	if cfgFilepath == "" && cfgInline == "" {
		err := errors.New("Required configuration filepath or inline configuration.")
		slog.Error(err)
		panic(err)
	}

	err := webserver.Initialize(cfgFilepath, cfgInline)
	if err != nil {
		slog.Error(err)
		panic(err)
	}

	handler := webserver.Routers()

	dataServerShutdownComplete := &sync.WaitGroup{}
	dataServerShutdownComplete.Add(1)
	dataServer := startHttpServer(":8080", handler, 1, dataServerShutdownComplete)

	// This function blocks the execution, once the signal is received the close up routine starts.
	catchKillSignal()
	log.Info("Received kill signal")

	log.Info("Closing webserver.")
	// Shutdown closes the server and subtracts 1 to the waitGroup, so execution can continue.
	if err := dataServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Error(fmt.Sprintf("apiServer Shutdown: %v", err))
	}
	dataServerShutdownComplete.Wait()

	log.Info("Closing Writers.")
	webserver.CloseWriters()

	log.Info("Shutdown complete.")
}

// This function starts the http server and returns the pointer to the server instance.
// With the pointer, the server can be shutdown later.
func startHttpServer(port string, h http.Handler, timeoutMultiplier int, wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{
		Addr:         port,
		Handler:      h,
		ReadTimeout:  1 * time.Duration(timeoutMultiplier) * time.Minute,
		WriteTimeout: 2 * time.Duration(timeoutMultiplier) * time.Minute,
	}

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err.Error())
		}
	}()

	return srv
}

// This function waits for a termination signal like Control + 'C' or kill command.
// It allows to have a graceful termination, which may be useful to finish writing data before the app is closed.
func catchKillSignal() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	<-sigint
}
