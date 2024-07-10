package organization

import (
	"errors"
	"fmt"
	"ospm/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetPolicyCheck(context *fiber.Ctx) (models.APIError, error) {
	if !strings.HasSuffix(context.Path(), "profile") && context.Query("list_all") == "true" {
		if ClientIPCanListAllOrganization(context.IP()) {
			return models.APIError{}, nil
		} else {
			return models.APIError{
				Error:   fiber.ErrBadRequest.Error(),
				Message: "the client is not permitted to list all organizations",
			}, errors.New("")
		}
	}

	return models.APIError{}, nil
}

func PatchPolicyCheck(context *fiber.Ctx) (models.APIError, error) {
	if ClientIPCanUndoOrganizationSoftDelete(context.IP()) {
		return models.APIError{}, nil
	} else {
		return models.APIError{
			Error:   fiber.ErrBadRequest.Error(),
			Message: "the client is not permitted to undo organization soft delete",
		}, errors.New("")
	}
}

func DeletePolicyCheck(context *fiber.Ctx) (models.APIError, error) {
	switch context.Query("mode") {
	case "soft":
		if ClientIPCanSoftDeleteOrganization(context.IP()) {
			return models.APIError{}, nil
		} else {
			return models.APIError{
				Error:   fiber.ErrBadRequest.Error(),
				Message: fmt.Sprintf("request from %s is not permitted to soft delete the organization", context.IP()),
			}, errors.New("")
		}
	case "hard":
		if ClientIPCanHardDeleteOrganization(context.IP()) {
			return models.APIError{}, nil
		} else {
			return models.APIError{
				Error:   fiber.ErrBadRequest.Error(),
				Message: fmt.Sprintf("request from %s is not permitted to hard delete the organization", context.IP()),
			}, errors.New("")
		}
	default:
		return models.APIError{
			Error:   fiber.ErrBadRequest.Error(),
			Message: "the deletion mode should be provided. valid values are: soft/hard",
		}, errors.New("")

	}
}
