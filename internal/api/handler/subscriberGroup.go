package handler

import (
	"encoding/json"
	"ospm/internal/models"
	"ospm/internal/service/organization"

	"github.com/gofiber/fiber/v2"
)

func GetOrganizationGroupList(context *fiber.Ctx) error {
	organizationID := context.Params("organization_id")

	organizationGroupList, err := organization.GetOrganizationGroupList(organizationID)
	if err != nil {
		errorMessage := models.APIError{
			Error:   err.Error(),
			Message: "failed to list the subscriber groups",
		}
		return context.Status(fiber.ErrInternalServerError.Code).JSON(errorMessage)
	} else if len(organizationGroupList) == 0 {
		return context.SendStatus(fiber.ErrNotFound.Code)
	}

	listInJson, err := json.Marshal(organizationGroupList)
	if err != nil {
		errorMessage := models.APIError{
			Error:   err.Error(),
			Message: "failed to load the subscriber groups",
		}
		return context.Status(fiber.ErrInternalServerError.Code).JSON(errorMessage)
	}

	return context.Status(200).Send(listInJson)
}
