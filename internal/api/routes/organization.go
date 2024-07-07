package routes

import (
	"ospm/internal/api/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupOrganizationRoutes(rg fiber.Router) {
	rg.Get("", handler.GetOrganizationList)
	rg.Get("/profile", handler.GetOrganizationProfile)
	rg.Post("", handler.AddNewOrganization)
	// rg.Patch("")
	// rg.Delete("")
}
