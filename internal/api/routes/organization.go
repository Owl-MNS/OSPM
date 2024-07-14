package routes

import (
	"ospm/internal/api/handler"
	"ospm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupOrganizationRoutes(rg fiber.Router) {
	rg.Get("", middleware.OrganizationPolicyCheck, handler.GetOrganizationList)
	rg.Get("/profile", middleware.OrganizationPolicyCheck, handler.GetOrganizationProfile)
	rg.Get("/:organization_id/subscriber-group", handler.GetSubscriberGroupList)
	rg.Post("", middleware.OrganizationPolicyCheck, handler.AddNewOrganization)
	rg.Delete("", middleware.OrganizationPolicyCheck, handler.DeleteOrganization)
	rg.Patch("/recover/profile", middleware.OrganizationPolicyCheck, handler.RecoverSoftDeletedOrganization)
}
