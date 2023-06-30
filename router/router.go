package router

import (
	"blog/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Router() {
	app := fiber.New()
	app.Post("/signup", handlers.Signup)
	app.Post("/login", handlers.Login)

	//start a server
	log.Println("Server starts in port 8000.....")
	app.Listen("localhost:8000")
}
