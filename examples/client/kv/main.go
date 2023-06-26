package main

func main() {

	// BOO: This still uses the old nats package APIs...
	// You would need to create the jetstream connection the old way.

	// nc, _ := nats.Connect(nats.DefaultURL)
	// js, _ := jetstream.New(nc)

	// js, _ := nc.JetStream()
	// defer nc.Close()

	// js.KeyValueStoreNames()
	// s.CreateKeyValue
}
