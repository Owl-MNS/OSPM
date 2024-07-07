package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	SetupOrganizationRoutes(app.Group("/organization"))

}
