package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jgcaceres97/go-auth-jwt/src/database"
	"github.com/jgcaceres97/go-auth-jwt/src/routes"
	"github.com/jgcaceres97/go-auth-jwt/src/settings"
)

func init() {
	err := settings.New()
	if err != nil {
		log.Panic(err)
	}

	database.Connect()
}

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{AllowCredentials: true}))

	routes.Setup(app)

	go func() {
		addr := fmt.Sprintf(":%s", *settings.Port)

		if err := app.Listen(addr); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("\nGracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	fmt.Println("\t- Closing db connection...")
	DB, _ := database.DB.DB()
	defer DB.Close()

	fmt.Println("\nServer was successfull shutdown.")
}
