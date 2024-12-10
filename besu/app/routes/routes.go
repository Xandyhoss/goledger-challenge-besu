package routes

import (
	"github.com/Xandyhoss/goledger-challenge-besu/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/set", controllers.ExecContractHandler)
	app.Get("/get", controllers.CallContractHandler)
	app.Get("/check", controllers.CheckContractHandler)
	app.Put("/sync", controllers.SyncContractHandler)
}
