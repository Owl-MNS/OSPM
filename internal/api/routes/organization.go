package routes

import (
	"ospm/internal/api/handler"
	"ospm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupOrganizationRoutes(rg fiber.Router) {

	// Organization => Subscriber Group Management
	rg.Get("/:organization-id/subscriber-group", handler.GetSubscriberGroupList)
	rg.Get("/:organization-id/subscriber-group/:subscriber-group-id", handler.GetSubscriberGroupDetail)
	rg.Post("/:organization-id/subscriber-group/:subscriber-group-id", handler.AddNewSubscriberGroup)
	rg.Delete("/:organization-id/subscriber-group/:subscriber-group-id", handler.DeleteSubscriberGroup)
	rg.Patch("/:organization-id/subscriber-group/:subscriber-group-id", handler.UpdateSubscriberGroup)

	// Organization Management
	rg.Get("", middleware.OrganizationPolicyCheck, handler.GetOrganizationList)
	rg.Get("/profile", middleware.OrganizationPolicyCheck, handler.GetOrganizationProfile)
	rg.Post("", middleware.OrganizationPolicyCheck, handler.AddNewOrganization)
	rg.Delete("", middleware.OrganizationPolicyCheck, handler.DeleteOrganization)
	rg.Patch("/recover/profile", middleware.OrganizationPolicyCheck, handler.RecoverSoftDeletedOrganization)
}
