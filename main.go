package main

import (
	"GpsTracker2/database"
	_ "GpsTracker2/docs"
	"GpsTracker2/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
	"log"
)

func main() {
	app := fiber.New() // It starts HTTP server

	log.Println("Starting GpsTracker2...")

	//database.InitRedis()
	database.ConnectMongo()
	database.ConnectMysql()
	database.ConnectElastic()

	routes.ProtectedUserRoutes(app)
	routes.ProtectedDeviceRoutes(app)
	routes.ProtectedLocationRoutes(app)

	app.Static("/swagger", "./docs")

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Handle websocket connection
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			// Echo the message back to the client
			if err = c.WriteMessage(mt, msg); err != nil {
				break
			}
		}
	}))

	// Serve Swagger UI that points to the YAML file
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "http://localhost:3000/swagger/swagger.yaml", // Point to your manual YAML file
	}))

	go func() {
		app.Listen(":3000") // It listens HTTP requests on port 3000
	}()

	go database.StartConsumer()

	select {}
}
