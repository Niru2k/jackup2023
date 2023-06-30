package router

import (
	//User-defined packages
	"blog/handlers"
	"blog/logs"

	//Third-party package
	"github.com/gofiber/fiber/v2"
)

func Router() {
	log := logs.Log()
	app := fiber.New()
	app.Post("/signup", handlers.Signup)
	app.Post("/login", handlers.Login)

	//start a server
	log.Info("Server starts in port 8000.....")
	app.Listen(":8000")
}
