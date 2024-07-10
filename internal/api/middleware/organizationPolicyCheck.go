package middleware

import (
	"ospm/internal/models"
	"ospm/internal/service/organization"

	"github.com/gofiber/fiber/v2"
)

// OrganizationPolicyCheck applies the policy checks based on the request method
func OrganizationPolicyCheck(context *fiber.Ctx) error {

	switch context.Method() {
	case "GET":
		apiError, err := organization.GetPolicyCheck(context)
		if err != nil {
			return context.Status(fiber.ErrBadRequest.Code).JSON(apiError)
		}
		return context.Next()

	case "DELETE":
		apiError, err := organization.DeletePolicyCheck(context)
		if err != nil {
			return context.Status(fiber.ErrBadRequest.Code).JSON(apiError)
		}
		return context.Next()

	case "POST":
		return context.Next()

	case "PATCH":
		apiError, err := organization.PatchPolicyCheck(context)
		if err != nil {
			return context.Status(fiber.ErrBadRequest.Code).JSON(apiError)
		}
		return context.Next()
	}

	return context.Status(fiber.ErrInternalServerError.Code).JSON(models.APIError{
		Error:   fiber.ErrBadRequest.Error(),
		Message: "unsupported request method",
	})
}
