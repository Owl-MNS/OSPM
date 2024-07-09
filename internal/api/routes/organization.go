package routes

import (
	"ospm/internal/api/handler"
	"ospm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupOrganizationRoutes(rg fiber.Router) {
	rg.Get("", middleware.OrganizationPolicyCheck, handler.GetOrganizationList)
	rg.Get("/profile", middleware.OrganizationPolicyCheck, handler.GetOrganizationProfile)
	rg.Post("", middleware.OrganizationPolicyCheck, handler.AddNewOrganization)
	rg.Delete("", middleware.OrganizationPolicyCheck, handler.DeleteOrganization)
	// rg.Patch("")
}
