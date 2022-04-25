package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/centrifugal/centrifuge"
)

const shutdownTimeout = 5 * time.Second

var host, port string

// handleLog logs messages from Centrifuge.
func handleLog(e centrifuge.LogEntry) {
	log.Printf("%s: %v", e.Message, e.Fields)
}

// setupSignalHandler registers for SIGTERM and SIGINT. A context is returned
// which is canceled on one of these signals. If a second signal is caught, the
// program is terminated with exit code 1.
func setupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)
	go func() {
		<-c
		cancel()

		// If signalled twice, immediately exit.
		<-c
		os.Exit(1)
	}()

	return ctx
}

// app encapsulates the Centrifuge and HTTP servers.
type app struct {
	bind string
	node *centrifuge.Node
	h    *http.Server
}

// run starts the Centrifuge and HTTP servers and runs them to completion.
func (a *app) run(ctx context.Context) error {
	// Start Centrifuge
	if err := a.node.Run(); err != nil {
		return fmt.Errorf("run centrifuge: %w", err)
	}
	log.Println("centrifuge running...")

	ch := make(chan error)
	defer close(ch)

	// Start HTTP server.
	go a.serveHTTP(ch)

	log.Printf("HTTP listening on %s\n", a.bind)

	// Wait on our context to be cancelled.
	<-ctx.Done()

	// Shutdown within 5 seconds.
	ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.node.Shutdown(ctxShutdown); err != nil {
		return err
	}
	if err := suppressServerClosed(a.h.Shutdown(ctxShutdown)); err != nil {
		return err
	}

	return nil
}

// suppressServerClosed hides the HTTP server closed error as it's expected.
func suppressServerClosed(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func main() {
	flag.StringVar(&host, "host", "", "IP on which to bind the server")
	flag.StringVar(&port, "port", "8888", "Port on which to bind the server")
	flag.Parse()

	a := app{
		bind: net.JoinHostPort(host, port),
	}

	if err := a.initCentrifuge(); err != nil {
		log.Fatalln(err)
	}

	if err := a.initHTTP(); err != nil {
		log.Fatalln(err)
	}

	// Start the application, shutting down on a signal.
	ctx := setupSignalHandler()
	log.Fatal(a.run(ctx))
}
