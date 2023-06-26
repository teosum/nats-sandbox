package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	fmt.Println("Starting consumer", start.String())
	defer func() {
		fmt.Println("Stopped - took", time.Since(start).String())
	}()

	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := jetstream.New(nc)

	// Create a stream - this call is idempotent
	s, _ := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})

	// Create durable consumer
	c, _ := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})

	// OPTION #1 - BATCH FETCH (continues to pull until batch size limit)
	// -----------------------------------------------------------------------------------------------
	// jetstream.MessageBatch interface contains two functions returning
	// 1. messages channel
	// 2. error string
	// msgs, _ := c.Fetch(100)
	// for msg := range msgs.Messages() {
	// 	_ = msg.Ack()
	// 	fmt.Printf("Received a JetStream message: %s\n", string(msg.Data()))
	// }

	// if msgs.Error() != nil {
	// 	fmt.Println("Error during Fetch(): ", msgs.Error())
	// }

	// OPTION #2 - CONSUME WITH CALLBACK
	// -----------------------------------------------------------------------------------------------
	cons, _ := c.Consume(func(msg jetstream.Msg) {
		_ = msg.Ack()
		dataBytes := msg.Data()
		fmt.Printf("Callback for message: %s\n", string(dataBytes))
	})
	defer cons.Stop()

	for range ctx.Done() {
		cons.Stop()
		return
	}

	// OPTION #3 - CONSUME
	// -----------------------------------------------------------------------------------------------
	// msgs, err := c.Messages()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for {
	// 	if ctx.Err() != nil {
	// 		msgs.Stop()
	// 		break
	// 	}

	// 	msg, err := msgs.Next()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	_ = msg.Ack()
	// 	fmt.Println(string(msg.Data()))
	// }
}
