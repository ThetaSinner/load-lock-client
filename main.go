package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

// Registration model
type registration struct {
	Id    string
	Group string
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Failed to ping redis server")
		os.Exit(1)
	}

	var registration = &registration{
		Id:    "abc-123",
		Group: "abc"}
	var msg, _ = json.Marshal(registration)

	subscription := client.Subscribe("load-lock:start:abc-123").Channel()

	client.LPush("load-lock:registration-queue", msg)

	fmt.Println("Listening on subscribed channels...")
	<-subscription

	fmt.Println("All done, this can now run")
}
