package routes

import (
	"ospm/internal/api/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupSubscriberGroupRoutes(rg fiber.Router) {

	rg.Get("/list/:organization_id", handler.GetSubscriberGroupList)
	rg.Get("/:subscriber_group_id", handler.GetSubscriberGroupDetail)
	rg.Post("/:organization_id", handler.AddNewSubscriberGroup)
	rg.Delete("/:subscriber_group_id", handler.DeleteSubscriberGroup)
	rg.Patch("/:subscriber-group-id", handler.UpdateSubscriberGroup)
}
