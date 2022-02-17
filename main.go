package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rama-kairi/fiber-api/config"
	"github.com/rama-kairi/fiber-api/database"
	"github.com/rama-kairi/fiber-api/routes"
)

func main() {
	database.ConnectDB()

	app := fiber.New(
		fiber.Config{
			Prefork:       false,
			AppName:       config.GetConfig().App.Name,
			StrictRouting: true,
		},
	)

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(config.GetConfig().App.Port))
}
