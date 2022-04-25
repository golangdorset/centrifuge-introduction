package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/centrifugal/centrifuge"
)

// initCentrifuge initialises the Centrifuge server.
func (a *app) initCentrifuge() error {
	var err error
	a.node, err = centrifuge.New(centrifuge.Config{
		LogLevel:   centrifuge.LogLevelInfo,
		LogHandler: handleLog,
	})
	if err != nil {
		return fmt.Errorf("initialise centrifuge: %w", err)
	}

	// Override default broker which does not use HistoryMetaTTL. This allows
	// the broker to automatically expire old messages.
	broker, err := centrifuge.NewMemoryBroker(a.node, centrifuge.MemoryBrokerConfig{
		HistoryMetaTTL: 2 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("new memory broker: %w", err)
	}
	a.node.SetBroker(broker)

	a.node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		return centrifuge.ConnectReply{}, nil
	})

	a.node.OnConnect(func(client *centrifuge.Client) {
		transport := client.Transport()
		log.Printf(
			"client connected via %s with protocol: %s", transport.Name(), transport.Protocol(),
		)

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			log.Printf("client subscribed to channel %s", e.Channel)

			cb(centrifuge.SubscribeReply{}, nil)
		})
	})

	return nil
}
