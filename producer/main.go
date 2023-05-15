package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	// "time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var redis_client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var ctx = context.Background()

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
	DateTime string `json:"date_time"`
}

func main() {

	// Setup app server
	app := fiber.New()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Hello 👋!")
	})

	app.Post("/create", func(c *fiber.Ctx) error {
		
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			panic(err)
		}

		payload, err := json.Marshal(user)

		if err != nil {
			panic(err)
		}
		
		if err := redis_client.Publish( ctx, "send-user-data", payload ).Err(); err != nil {
			panic(err)
		}

		return c.SendStatus(200)
	})

	// Listen to a port
	if error := app.Listen(":6969"); error != nil {
		log.Panic(error)
	}

	fmt.Println("Running cleanup tasks...")
}
