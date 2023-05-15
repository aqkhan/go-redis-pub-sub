package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
	DateTime string `json:"date_time"`
}

type Task struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

var ctx = context.Background()

var redis_client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	fmt.Println("Hello 👋🏻 from consumer")

	// Create subscriber
	subscriber := redis_client.Subscribe(ctx, "send-user-data")

	user := User{}

	// Wait for message
	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal([]byte(msg.Payload), &user); err != nil {
			panic(err)
		}

		fmt.Println("Received data from: " + msg.Channel + " channel 🎉")
		fmt.Printf("%+v\n", user)
	}
}