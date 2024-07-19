package subscriberGroup

import (
	"fmt"
	"ospm/internal/models"
	"ospm/internal/repository/database/cockroachdb"
	"ospm/internal/service/logger"
	"reflect"
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
func NewByAPI(newSubscriberGroup models.CreateUpdateSubscriberGroupAPI) (string, error) {

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

func Update(newSubscriberGroupDetails models.CreateUpdateSubscriberGroupAPI, subscriberGroupID string) error {

	currentSubscriberGroup, err := Detail(subscriberGroupID)
	if err != nil {
		errorMessage := fmt.Sprintf("failed to load the subscriber group settings for the given id: %s, error: %+v",
			subscriberGroupID, err.Error())
		logger.OSPMLogger.Errorln(errorMessage)
		return err
	}

	updateTX := cockroachdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			errorMessage := fmt.Sprintf("failed to update the subscriber group for the given id: %s, error: %+v",
				subscriberGroupID, r)
			logger.OSPMLogger.Errorln(errorMessage)
			updateTX.Rollback()
		}
	}()

	if currentSubscriberGroup.Name != newSubscriberGroupDetails.Name {
		err = updateTX.Unscoped().Model(&models.SubscriberGroup{}).Where("id = ?", subscriberGroupID).Update("subscriber_group_name",
			newSubscriberGroupDetails.Name).Error
		if err != nil {
			errorMessage := fmt.Sprintf("failed to update the subscriber group name at update step for the given id: %s, error: %+v",
				subscriberGroupID, err.Error())
			logger.OSPMLogger.Errorln(errorMessage)
			updateTX.Rollback()
			return err
		}
	}

	if currentSubscriberGroup.Description != newSubscriberGroupDetails.Description {
		err = updateTX.Unscoped().Model(&models.SubscriberGroup{}).
			Where("id = ?", subscriberGroupID).
			Update("description", newSubscriberGroupDetails.Description).Error
		if err != nil {
			errorMessage := fmt.Sprintf("failed to update the subscriber group name at update step for the given id: %s, error: %+v",
				subscriberGroupID, err.Error())
			logger.OSPMLogger.Errorln(errorMessage)
			updateTX.Rollback()
			return err
		}
	}

	for permCategory, permSettings := range newSubscriberGroupDetails.Permissions {
		for _, perm := range permSettings {
			updateCandidatePerm := models.Permission{
				PermissionCategory: permCategory,
				PermissionName:     perm.PermissionName,
				PermissionValue:    perm.PermissionValue,
			}

			err = nil
			switch action(updateCandidatePerm, currentSubscriberGroup.Permissions) {
			case "update":
				err = updateTX.Unscoped().Model(&models.Permission{}).
					Where("subscriber_group_id = ? AND permission_name = ? AND permission_category = ?",
						subscriberGroupID,
						updateCandidatePerm.PermissionName,
						updateCandidatePerm.PermissionCategory).
					Update("permission_value", updateCandidatePerm.PermissionValue).Error
			case "create":
				updateCandidatePerm.SubscriberGroupID = subscriberGroupID
				err = updateTX.Create(&updateCandidatePerm).Error
			default:
				logger.OSPMLogger.Infoln("permission already exists. ignored. permission: %+v, Subscriber Group ID: %s", updateCandidatePerm, subscriberGroupID)
			}

			if err != nil {
				errorMessage := fmt.Sprintf(
					"failed to update the subscriber group permission at update step for the given id: %s, permission to update: %+v error: %+v",
					subscriberGroupID, updateCandidatePerm, err.Error())
				logger.OSPMLogger.Errorln(errorMessage)
				updateTX.Rollback()
				return err
			}
		}
	}

	err = updateTX.Commit().Error
	if err != nil {
		errorMessage := fmt.Sprintf("failed to update the subscriber group name at commit step for the given id: %s, error: %+v",
			subscriberGroupID, err.Error())
		logger.OSPMLogger.Errorln(errorMessage)
		updateTX.Rollback()
		return err
	}

	return nil
}

// action gets a permission config and checks it among a list of permissions and returns:
// ignore: if the permission list has the exact permission
// update: if the permission list has a permission with the same settings of givem permission but value
// create: if the permission list has not permission equal to the givem permission
func action(permToCheck models.Permission, PermsList []models.Permission) string {
	for _, perm := range PermsList {
		if reflect.DeepEqual(perm, permToCheck) {
			return "ignore"
		} else if perm.PermissionName == permToCheck.PermissionName &&
			perm.PermissionCategory == permToCheck.PermissionCategory &&
			perm.PermissionValue != permToCheck.PermissionValue {
			return "update"
		}
	}

	return "create"
}
