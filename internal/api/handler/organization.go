package handler

import (
	"encoding/json"
	"fmt"
	"ospm/internal/models"
	"ospm/internal/service/organization"

	"github.com/gofiber/fiber/v2"
)

// @Summary 	List all organizations
// @Description Retrieves a list of all organizations available in the system. each record is summarized
// @Tags 		Organization
// @Produce 	json
// @Success 	200 {array} models.OrganizationShortInfo "Successful Response"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/organizations [get]
func GetOrganizationList(context *fiber.Ctx) error {
	organizationList, err := organization.List()
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(models.APIError{
			Error:   fiber.ErrInternalServerError.Error(),
			Message: err.Error(),
		})
	}

	if len(organizationList) == 0 {
		return context.Status(fiber.ErrNotFound.Code).JSON(models.APIError{
			Error:   fiber.ErrNotFound.Error(),
			Message: fiber.ErrNotFound.Message,
		})
	}

	listInJson, err := json.Marshal(organizationList)
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(models.APIError{
			Error:   fiber.ErrInternalServerError.Error(),
			Message: err.Error(),
		})
	}

	return context.Status(200).Send(listInJson)
}

// @Summary 	Get organization profile by name or ID
// @Description Retrieves detailed information about a specific organization identified by its name or ID.
// @Tags 		Organizations
// @Accept 		json
// @Produce 	json
// @Param 		name query string false "Organization Name" @in query
// @Param 		id query string false "Organization ID" @in query
// @Success 	200 {object} models.OrganizationResponse "Successful Response"
// @Failure 	404 {object} models.APIError "Organization Not Found"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/organizations/profile [get]
func GetOrganizationProfile(context *fiber.Ctx) error {
	organizationName := context.Query("name")
	organizationID := context.Query("id")

	if organizationID == "" && organizationName == "" {
		return context.Status(fiber.StatusBadRequest).JSON(models.APIError{
			Error:   fiber.ErrBadRequest.Error(),
			Message: "Either organization ID or name must be provided",
		})
	}

	organizationDetails, err := organization.Details(organizationName, organizationID)
	if err != nil {
		status := fiber.StatusInternalServerError
		message := fiber.ErrInternalServerError.Message

		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
			message = err.Error()
		}

		return context.Status(status).JSON(models.APIError{
			Error:   err.Error(),
			Message: message,
		})
	}

	detailsInJson, err := json.Marshal(organization.Clean(&organizationDetails))
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(models.APIError{
			Error:   fiber.ErrInternalServerError.Error(),
			Message: err.Error(),
		})
	}

	return context.Status(fiber.StatusOK).Send(detailsInJson)
}

// @Summary 	Add a new organization
// @Description Adds a new organization to the system. The request body must contain the organization details.
// @Tags 		Organizations
// @Accept 		json
// @Produce 	json
// @Param 		body body models.Organization true "Organization details"
// @Success 	201 {object} map[string]string "Organization successfully added"
// @Failure 	400 {object} models.APIError "Bad Request"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/organization [post]
func AddNewOrganization(context *fiber.Ctx) error {
	newOrganization := models.Organization{}
	err := context.BodyParser(&newOrganization)
	if err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(models.APIError{
			Error:   fiber.ErrBadRequest.Error(),
			Message: "failed to pars the provided information, error:" + err.Error(),
		})
	}

	newOrganizationID, err := organization.New(newOrganization)
	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(models.APIError{
			Error:   fiber.ErrInternalServerError.Error(),
			Message: "failed to add new organization, error:" + err.Error(),
		})
	}

	responseMessage := map[string]string{
		"message":             fmt.Sprintf("organization %s successfulle added", newOrganization.Details.Name),
		"new_organization_id": newOrganizationID,
	}

	return context.Status(201).JSON(responseMessage)
}

// @Summary 	Delete an organization
// @Description Deletes an existing organization from the system. The organization can be specified by either its ID or name.
// @Tags 		Organizations
// @Accept 		json
// @Produce 	json
// @Param 		id query string false "Organization ID"
// @Param 		name query string false "Organization Name"
// @Success 	200 {object} map[string]string "Organization successfully deleted"
// @Failure 	400 {object} models.APIError "Bad Request"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/organizations/delete [delete]
func DeleteOrganization(context *fiber.Ctx) error {
	organizationName := context.Query("name")
	organizationID := context.Query("id")

	if organizationID == "" && organizationName == "" {
		return context.Status(fiber.StatusBadRequest).JSON(models.APIError{
			Error:   fiber.ErrBadRequest.Error(),
			Message: "Either organization ID or name must be provided",
		})
	}

	if err := organization.Delete(organizationID, organizationName); err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(models.APIError{
			Error:   err.Error(),
			Message: "failed to delete the organization",
		})
	}

	responseMessage := map[string]string{
		"message":                "organization successfulle deleted",
		"organization_to_delete": fmt.Sprintf("%s %s", organizationID, organizationName),
	}

	return context.Status(200).JSON(responseMessage)
}
