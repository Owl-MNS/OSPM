package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	SetupAPIDocs(app.Group("/apidoc"))
	SetupOrganizationRoutes(app.Group("/organization"))

}
