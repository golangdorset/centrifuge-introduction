package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/centrifugal/centrifuge-go"
)

// eventHandler implements centrifuge.PublishHandler
type eventHandler struct{}

func (h *eventHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	log.Printf("Someone says via channel %s: %s", sub.Channel(), string(e.Data))
}

var (
	userID string
	host   string
	port   string
)

// subscribe creates a subscription to a channel.
func subscribe(client *centrifuge.Client, channel string, handler centrifuge.PublishHandler) error {
	sub, err := client.NewSubscription(channel)
	if err != nil {
		return fmt.Errorf("new subscription: %w", err)
	}
	sub.OnPublish(handler)

	if err = sub.Subscribe(); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	return nil
}

func main() {
	flag.StringVar(&userID, "user", "123", "user ID")
	flag.StringVar(&host, "host", "localhost", "host of the server")
	flag.StringVar(&port, "port", "8888", "port of the server")
	flag.Parse()

	url := fmt.Sprintf(
		"ws://%s/v1/connection/websocket", net.JoinHostPort(host, port),
	)

	// Create a Centrifuge client.
	client := centrifuge.NewJsonClient(url, centrifuge.DefaultConfig())
	defer client.Close()

	handler := &eventHandler{}

	if err := subscribe(client, "broadcast", handler); err != nil {
		log.Fatalln(err)
	}
	if err := subscribe(client, userID, handler); err != nil {
		log.Fatalln(err)
	}

	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)
	<-c
}
