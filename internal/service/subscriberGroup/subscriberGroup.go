package subscriberGroup

import (
	"errors"
	"fmt"
	"ospm/internal/models"
	"ospm/internal/repository/database/cockroachdb"
	"ospm/internal/service/logger"
)

// GetSubscriberGroupList get the organization id and returns all groups within the given organiztion
// In the listed group, soft deleted groups are excluded!
func List(organizationsID string) ([]models.SubscriberGroupMinimal, error) {
	var groupList []models.SubscriberGroupMinimal

	err := cockroachdb.DB.
		Model(&models.SubscriberGroup{}).
		Select("id,name").
		Where("organization_id =  ?", organizationsID).
		Find(&groupList).Error
	if err != nil {
		errorMessage := fmt.Sprintf("failed to load the list of subscriber group for organization id %s, error: %s", organizationsID, err.Error())
		logger.OSPMLogger.Errorln(errorMessage)
		return nil, err
	}

	return groupList, nil
}

func Detail(subscriberGroupID string) (models.SubscriberGroup, error) {
	var subscriberGroupDetail models.SubscriberGroup
	err := cockroachdb.DB.Preload("Permissions").First(&subscriberGroupDetail, "id = ? ", subscriberGroupID).Error
	if err != nil {
		errorMessage := fmt.Sprintf("failed to load details of given group id %s, error: %+v", subscriberGroupID, err)
		logger.OSPMLogger.Errorln(errorMessage)
		return models.SubscriberGroup{}, err
	}

	return subscriberGroupDetail, nil
}

func Delete(subscriberGroupID string) error {
	deletetionTX := cockroachdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			deletetionTX.Rollback()
			errorMessage := fmt.Sprintf(
				"failed to delete the given group id %s, error: %+v",
				subscriberGroupID, r)
			logger.OSPMLogger.Errorln(errorMessage)
		}
	}()

	err := deletetionTX.Unscoped().Where("id = ?", subscriberGroupID).Delete(&models.SubscriberGroup{}).Error
	if err != nil {
		deletetionTX.Rollback()
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete group step, error: %+v",
			subscriberGroupID, err.Error())
		logger.OSPMLogger.Errorln(errorMessage)
		return err
	}

	err = deletetionTX.Unscoped().Where("subscriber_group_id = ?", subscriberGroupID).Delete(&models.Permission{}).Error
	if err != nil {
		deletetionTX.Rollback()
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete permission set step, error: %+v",
			subscriberGroupID, err.Error())
		logger.OSPMLogger.Errorln(errorMessage)
		return err
	}

	err = deletetionTX.Commit().Error
	if err != nil {
		deletetionTX.Rollback()
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at apply step, error: %+v",
			subscriberGroupID, err.Error())
		logger.OSPMLogger.Errorln(errorMessage)
		return err
	}

	logger.OSPMLogger.Infoln("subscriber group id %s successfully deleted", subscriberGroupID)

	return nil
}

// New gets the new subscriber group configuration based on standard subscriber group model and adds it
func New(newSubscriberGroup models.SubscriberGroup) (string, error) {
	createTX := cockroachdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			createTX.Rollback()
			errorMessage := fmt.Sprintf(
				"failed to add the new subscriber group  %s at apply step, error: %+v",
				newSubscriberGroup.Name, r)
			logger.OSPMLogger.Errorln(errorMessage)
		}
	}()

	err := createTX.Create(&newSubscriberGroup).Error
	if err != nil {
		createTX.Rollback()
		errorMessage := fmt.Sprintf(
			"failed to add the new subscriber group  %s at apply step, error: %+v",
			newSubscriberGroup.Name, err)
		logger.OSPMLogger.Errorln(errorMessage)
		return "-1", err
	}

	err = createTX.Commit().Error
	if err != nil {
		createTX.Rollback()
		errorMessage := fmt.Sprintf(
			"failed to add the new subscriber group  %s at apply step, error: %+v",
			newSubscriberGroup.Name, err)
		logger.OSPMLogger.Errorln(errorMessage)
		return "-1", err
	}

	logger.OSPMLogger.Infoln("subscriber group %s successfully added. id: %s", newSubscriberGroup.Name, newSubscriberGroup.ID)

	return newSubscriberGroup.ID, nil
}

// NewByAPI gets the new subscriber configurations based on a modified version of the standard
// model which has been used to ease the AddNewSubscriberGroup API Call
func NewByAPI(newSubscriberGroup models.SubscriberGroupAPI) (string, error) {

	//converting modified subscriber group used by API client to the standard version
	standardSubscriberGroup := models.SubscriberGroup{}
	standardSubscriberGroup.Absorb(newSubscriberGroup)

	createTX := cockroachdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			createTX.Rollback()
			errorMessage := fmt.Sprintf(
				"failed to add the new subscriber group  %s at apply step, error: %+v",
				newSubscriberGroup.Name, r)
			logger.OSPMLogger.Errorln(errorMessage)
		}
	}()

	err := createTX.Create(&standardSubscriberGroup).Error
	if err != nil {
		createTX.Rollback()
		errorMessage := fmt.Sprintf(
			"failed to add the new subscriber group  %s at apply step, error: %+v",
			newSubscriberGroup.Name, err)
		logger.OSPMLogger.Errorln(errorMessage)
		return "-1", err
	}

	err = createTX.Commit().Error
	if err != nil {
		createTX.Rollback()
		errorMessage := fmt.Sprintf(
			"failed to add the new subscriber group  %s at apply step, error: %+v",
			newSubscriberGroup.Name, err)
		logger.OSPMLogger.Errorln(errorMessage)
		return "-1", err
	}

	logger.OSPMLogger.Infoln("subscriber group %s successfully added. id: %s", standardSubscriberGroup.Name, standardSubscriberGroup.ID)

	return standardSubscriberGroup.ID, nil
}

func Update(newSubscriberGroupDetails models.SubscriberGroup, subscriberGroupID string) error {
	var oldSubscriberGroupDetail models.SubscriberGroup
	err := cockroachdb.DB.Preload("Permissions").
		Preload("Permissions.OrganizationLevelPerms").
		Preload("Permissions.AccessLevelPerms").
		Preload("Permissions.SubscriberLevelPerms").
		Preload("Permissions.PaymentLevelPerms").
		Preload("Permissions.ReportLevelPerms").
		Preload("Permissions.SubscriberGroupID").First(&oldSubscriberGroupDetail).Where("id = ? ", subscriberGroupID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to find the given group id %s to update, error: %+v",
			subscriberGroupID, err.Error())
		logger.OSPMLogger.Errorln(errorMessage)
		return err
	}

	if oldSubscriberGroupDetail.Name != newSubscriberGroupDetails.Name && newSubscriberGroupDetails.Name != "" {
		err = cockroachdb.DB.Model(&models.SubscriberGroup{}).
			Update("subscriber_group_name", newSubscriberGroupDetails.Name).
			Where("id = ?", subscriberGroupID).Error
		if err != nil {
			errorMessage := fmt.Sprintf(
				"failed to update the given group id %s, error: %+v",
				subscriberGroupID, err.Error())
			logger.OSPMLogger.Errorln(errorMessage)
			return err
		}
	}

	if oldSubscriberGroupDetail.Description != newSubscriberGroupDetails.Description && newSubscriberGroupDetails.Description != "" {
		err = cockroachdb.DB.Model(&models.SubscriberGroup{}).
			Update("subscriber_group_description", newSubscriberGroupDetails.Description).
			Where("id = ?", subscriberGroupID).Error
		if err != nil {
			errorMessage := fmt.Sprintf(
				"failed to update the given group id %s, error: %+v",
				subscriberGroupID, err.Error())
			logger.OSPMLogger.Errorln(errorMessage)
			return errors.New(errorMessage)
		}
	}

	// update perms should be added here

	logger.OSPMLogger.Infoln("subscriber group %s successfully updated. id: %s", newSubscriberGroupDetails.Name, subscriberGroupID)

	return nil
}
