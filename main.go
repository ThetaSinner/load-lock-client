package main

import (
	"encoding/json"
	"flag"
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
	var idFlag string
	flag.StringVar(&idFlag, "id", "", "The id that the client will use")

	var groupFlag string
	flag.StringVar(&groupFlag, "group", "", "The group that the client will use")

	isRegisterPtr := flag.Bool("register", false, "Whether to use this client session to register")

	flag.Parse()

	if idFlag == "" {
		panic("Missing --id flag")
	}

	if groupFlag == "" {
		panic("Missing --group flag")
	}

	if *isRegisterPtr {
		fmt.Printf("Will perform a registration. [id=%s], [group=%s]\n", idFlag, groupFlag)
		runRegistration(idFlag, groupFlag)
		return
	}

	fmt.Println("You didn't specify a valid command")
}

func createClient() *redis.Client {
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

	return client
}

func runRegistration(idFlag string, groupFlag string) {
	client := createClient()

	var registration = &registration{
		Id:    idFlag,
		Group: groupFlag}
	var msg, _ = json.Marshal(registration)

	subChannel := fmt.Sprintf("load-lock:start:%s", idFlag)
	subscription := client.Subscribe(subChannel).Channel()

	client.LPush("load-lock:registration-queue", msg)

	fmt.Println("Listening on subscribed channels...")
	<-subscription

	fmt.Println("All done, this can now run")
}
