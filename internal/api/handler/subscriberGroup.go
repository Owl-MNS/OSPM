package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"ospm/internal/models"
	"ospm/internal/service/logger"
	"ospm/internal/service/subscriberGroup"

	// This line is being used by swagger auto-documenting
	_ "ospm/docs/api"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// @Summary 	List all Subscriber Groups
// @Description Returns a list of all subscriber groups within an organization
// @Tags 		Organization
// @Accept  	json
// @Produce  	json
// @Param 		organization_id path int true "Organization ID"
// @Success 	200 {array} models.SubscriberGroupMinimal "Successful response"
// @Failure 	404 {object} models.APIError "Not Found"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/subscriber-group/list/{organization_id} [get]
func GetSubscriberGroupList(context *fiber.Ctx) error {
	organizationID := context.Params("organization_id")

	organizationGroupList, err := subscriberGroup.List(organizationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMessage := models.APIError{
				Error:   err.Error(),
				Message: "failed to load the subscriber group list",
			}
			return context.Status(fiber.ErrNotFound.Code).JSON(errorMessage)
		} else {
			errorMessage := models.APIError{
				Error:   err.Error(),
				Message: "failed to list the subscriber group list",
			}
			return context.Status(fiber.ErrInternalServerError.Code).JSON(errorMessage)
		}
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

// @Summary 	Get Subscriber Group Detail
// @Description Retrieves the details of a specific subscriber group within an organization
// @Tags 		Organization
// @Accept  	json
// @Produce  	json
// @Param 		organization_id path int true "Organization ID"
// @Param 		subscriber_group_id path int true "Subscriber Group ID"
// @Success 	200 {object} models.SubscriberGroupAPI "Successful response"
// @Failure 	404 {object} models.APIError "Not Found"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/subscriber_group/{subscriber_group_id} [get]
func GetSubscriberGroupDetail(context *fiber.Ctx) error {
	subscriberGroupID := context.Params("subscriber_group_id")

	groupDetail, err := subscriberGroup.Detail(subscriberGroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMessage := models.APIError{
				Error:   err.Error(),
				Message: "failed to load the subscriber group details",
			}
			return context.Status(fiber.ErrNotFound.Code).JSON(errorMessage)
		} else {
			errorMessage := models.APIError{
				Error:   err.Error(),
				Message: "failed to load the subscriber group details",
			}
			return context.Status(fiber.ErrInternalServerError.Code).JSON(errorMessage)
		}
	}

	jsonDetails, err := json.Marshal(groupDetail.Beautify())
	if err != nil {
		errorMessage := models.APIError{
			Error:   err.Error(),
			Message: "failed to load the subscriber group details",
		}
		return context.Status(fiber.ErrInternalServerError.Code).JSON(errorMessage)
	}

	return context.Status(200).Send(jsonDetails)
}

// @Summary 	Add New Subscriber Group
// @Description Adds a new subscriber group within an organization
// @Tags 		Organization
// @Accept  	json
// @Produce  	json
// @Param 		organization_id path int true "Subscriber Group ID"
// @Param 		body body models.CreateUpdateSubscriberGroupAPI true "Subscriber Group Details"
// @Success 	201 {object} models.SubscriberGroupCreateResponse "Successfully added new subscriber group"
// @Failure 	400 {object} models.APIError "Bad Request"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/subscriber_group/{organization_id} [post]
func AddNewSubscriberGroup(context *fiber.Ctx) error {
	var newSubscriberGroup models.CreateUpdateSubscriberGroupAPI

	err := context.BodyParser(&newSubscriberGroup)
	if err != nil {
		errorMessage := models.APIError{
			Error:   err.Error(),
			Message: "failed to process the request",
		}
		logger.OSPMLogger.Errorln(
			fmt.Sprintf(
				"failed to process request. Path: %s, client ip: %s, error: %+v",
				context.Path(), context.IP(), err))
		return context.Status(fiber.StatusBadRequest).JSON(errorMessage)
	}

	newSubscriberGroup.OrganizationID = context.Params("organization_id")
	id, err := subscriberGroup.NewByAPI(newSubscriberGroup)
	if err != nil {
		responseCode := 500
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			responseCode = fiber.ErrBadRequest.Code
		}
		errorMessage := models.APIError{
			Error:   err.Error(),
			Message: "failed to process the request",
		}
		logger.OSPMLogger.Errorln(
			fmt.Sprintf(
				"failed to process request. Path: %s, client ip: %s, error: %+v",
				context.Path(), context.IP(), err))
		return context.Status(responseCode).JSON(errorMessage)
	}

	response := models.SubscriberGroupCreateResponse{
		Message: "new subscriber group successfully added",
		Name:    newSubscriberGroup.Name,
		Id:      id,
	}

	return context.Status(201).JSON(response)
}

// @Summary 	Update Subscriber Group
// @Description Updates the settings of a specific subscriber group within an organization by its ID
// @Tags 		Organization
// @Accept  	json
// @Produce  	json
// @Param 		subscriber_group_id path int true "Subscriber Group ID"
// @Param 		body body models.CreateUpdateSubscriberGroupAPI true "Subscriber Group Settings"
// @Success 	200 "No Content"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/subscriber-group/{subscriber_group_id} [patch]
func UpdateSubscriberGroup(context *fiber.Ctx) error {
	var newSubscriberGroupSettings models.CreateUpdateSubscriberGroupAPI
	var responseCode int

	subscriberGroupID := context.Params("subscriber_group_id")

	err := context.BodyParser(&newSubscriberGroupSettings)
	if err != nil {
		errorMessage := models.APIError{
			Error:   err.Error(),
			Message: "failed to process the request",
		}
		logger.OSPMLogger.Errorln(
			fmt.Sprintf(
				"failed to process request. Path: %s, client ip: %s, error: %+v",
				context.Path(), context.IP(), err))
		return context.Status(fiber.StatusBadRequest).JSON(errorMessage)
	}

	err = subscriberGroup.Update(newSubscriberGroupSettings, subscriberGroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responseCode = fiber.ErrNotFound.Code
		} else {
			responseCode = fiber.ErrInternalServerError.Code
		}
		errorMessage := models.APIError{
			Error:   err.Error(),
			Message: "failed to process the request",
		}
		logger.OSPMLogger.Errorln(
			fmt.Sprintf(
				"failed to process request. Path: %s, client ip: %s, error: %+v",
				context.Path(), context.IP(), err))
		return context.Status(responseCode).JSON(errorMessage)
	}

	return context.SendStatus(200)

}

// @Summary 	Delete a Subscriber Group
// @Description Deletes a specific subscriber group within an organization by its ID
// @Tags 		Organization
// @Accept  	json
// @Produce  	json
// @Param 		subscriber-group-id path int true "Subscriber Group ID"
// @Param 		organization-id path int true "Subscriber Group ID"
// @Success 	204 "No Content"
// @Failure 	404 {object} models.APIError "Not Found"
// @Failure 	500 {object} models.APIError "Internal Server Error"
// @Router 		/subscriber-group/{subscriber-group-id} [delete]
func DeleteSubscriberGroup(context *fiber.Ctx) error {
	var responseCode int
	subscriberGroupID := context.Params("subscriber_group_id")

	err := subscriberGroup.Delete(subscriberGroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responseCode = fiber.ErrNotFound.Code
		} else {
			responseCode = fiber.ErrInternalServerError.Code
		}
		errorMessage := models.APIError{
			Error:   err.Error(),
			Message: "failed to process the request",
		}
		logger.OSPMLogger.Errorln(
			fmt.Sprintf(
				"failed to process request. Path: %s, client ip: %s, error: %+v",
				context.Path(), context.IP(), err))
		return context.Status(responseCode).JSON(errorMessage)
	}

	return context.SendStatus(204)
}
