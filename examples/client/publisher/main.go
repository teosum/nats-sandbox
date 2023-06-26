package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	// documentation for simplified JetStream client
	// https://github.com/nats-io/nats.go/blob/main/jetstream/README.md

	timeout := time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := jetstream.New(nc)

	// Create a stream - this call is idempotent
	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})

	info, err := s.Info(ctx)
	if err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(info, "", "  ")
	fmt.Printf("Created stream: %s\n\n", j)
	fmt.Printf("Publishing messages for %s\n", timeout.String())

	// Publish some messages
	t := time.NewTicker(time.Second)
	defer t.Stop()

	subject := "ORDERS.new"
	counter := 0

	for {
		select {
		case <-t.C:
			counter++
			msg := fmt.Sprintf("message %d", counter)

			js.Publish(ctx, subject, []byte(msg))
			fmt.Printf("Published: %s\n", msg)
		case <-ctx.Done():
			fmt.Println("context canceled")
			return
		}
	}
}
