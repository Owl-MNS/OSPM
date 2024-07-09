package middleware

import (
	"fmt"
	"ospm/internal/models"
	"ospm/internal/service/organization"

	"github.com/gofiber/fiber/v2"
)

func OrganizationPolicyCheck(context *fiber.Ctx) error {

	if context.Method() == "DELETE" {
		switch context.Query("mode") {
		case "soft":
			if organization.ClientIPCanSoftDeleteOrganization(context.IP()) {
				return context.Next()
			} else {
				return context.Status(fiber.ErrBadRequest.Code).JSON(models.APIError{
					Error:   fiber.ErrBadRequest.Error(),
					Message: fmt.Sprintf("request from %s is not permitted to soft delete the organization", context.IP()),
				})
			}
		case "hard":
			if organization.ClientIPCanHardDeleteOrganization(context.IP()) {
				return context.Next()
			} else {
				return context.Status(fiber.ErrBadRequest.Code).JSON(models.APIError{
					Error:   fiber.ErrBadRequest.Error(),
					Message: fmt.Sprintf("request from %s is not permitted to hard delete the organization", context.IP()),
				})
			}
		default:
			return context.Status(fiber.ErrBadRequest.Code).JSON(models.APIError{
				Error:   fiber.ErrBadRequest.Error(),
				Message: "the deletion mode should be provided. valid values are: soft/hard",
			})

		}
	}

	return context.Next()
}
