package main

import (
	"log"
	"os"

	"github.com/Xandyhoss/goledger-challenge-besu/app/db"
	"github.com/Xandyhoss/goledger-challenge-besu/app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.StartContractDB()

	app := fiber.New()

	routes.SetupRoutes(app)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"message": "Server is running OK!",
		})
	})

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
	log.Println("Servidor iniciado na porta " + os.Getenv("PORT"))
}
