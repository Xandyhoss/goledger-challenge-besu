package controllers

import (
	"github.com/Xandyhoss/goledger-challenge-besu/app/services"
	"github.com/gofiber/fiber/v2"
)

func ExecContractHandler(c *fiber.Ctx) error {
	type RequestBody struct {
		Value uint `json:"value"`
	}

	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	result, err := services.ExecContract(req.Value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": result,
	})
}

func CallContractHandler(c *fiber.Ctx) error {
	result, err := services.CallContract()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": result,
	})
}

func CheckContractHandler(c *fiber.Ctx) error {
	result := services.CheckContract()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": result,
	})
}

func SyncContractHandler(c *fiber.Ctx) error {
	services.SyncContract()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "contract synced",
	})
}
