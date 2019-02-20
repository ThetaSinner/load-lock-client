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
	ID    string
	Group string
}

func main() {
	var idFlag string
	flag.StringVar(&idFlag, "id", "", "The id that the client will use")

	var groupFlag string
	flag.StringVar(&groupFlag, "group", "", "The group that the client will use")

	isRegisterPtr := flag.Bool("register", false, "Whether to use this client session to register")

	isReleasePtr := flag.Bool("release", false, "Whether to use this client session to release")

	flag.Parse()

	if idFlag == "" {
		panic("Missing --id flag")
	}

	if *isRegisterPtr && groupFlag == "" {
		panic("Missing --group flag")
	}

	if *isRegisterPtr && *isReleasePtr {
		panic("Cannot register and release in the same session!")
	}

	if *isRegisterPtr {
		fmt.Printf("Will perform a registration. [id=%s], [group=%s]\n", idFlag, groupFlag)
		runRegistration(idFlag, groupFlag)
		return
	}

	if *isReleasePtr {
		fmt.Printf("Will perform a release. [id=%s]\n", idFlag)
		runRelease(idFlag)
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
		ID:    idFlag,
		Group: groupFlag}
	var msg, _ = json.Marshal(registration)

	subChannel := fmt.Sprintf("load-lock:start:%s", idFlag)
	subscription := client.Subscribe(subChannel).Channel()

	client.LPush("load-lock:registration-queue", msg)

	fmt.Println("Listening on subscribed channels...")
	<-subscription

	fmt.Println("All done, this can now run")
}

func runRelease(idFlag string) {
	client := createClient()

	client.LPush("load-lock:release-queue", idFlag)

	fmt.Println("Release notified. You can go about your day knowing you've done a good thing!")
}
