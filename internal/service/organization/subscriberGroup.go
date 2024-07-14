package organization

import (
	"errors"
	"fmt"
	"ospm/internal/models"
	"ospm/internal/repository/database/cockroachdb"
	OSPMLogger "ospm/internal/service/log"
)

// GetSubscriberGroupList get the organization id and returns all groups within the given organiztion
// In the listed group, soft deleted groups are excluded!
func GetSubscriberGroupList(organizationsID string) ([]models.SubscriberGroupMinimal, error) {
	var groupList []models.SubscriberGroupMinimal

	err := cockroachdb.DB.
		Model(&models.SubscriberGroup{}).
		Select("id,name").
		Where("organization_id =  ?", organizationsID).
		Find(&groupList).Error
	if err != nil {
		errorMessage := fmt.Sprintf("failed to load the list of subscriber group for organization id %s, error: %s", organizationsID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return groupList, nil
}

func GetSubscriberGroupDetail(subscriberGroupID string) (models.SubscriberGroup, error) {
	var subscriberGroupDetail models.SubscriberGroup
	err := cockroachdb.DB.Preload("Permissions").
		Preload("Permissions.OrganizationLevelPerms").
		Preload("Permissions.AccessLevelPerms").
		Preload("Permissions.SubscriberLevelPerms").
		Preload("Permissions.PaymentLevelPerms").
		Preload("Permissions.ReportLevelPerms").
		Preload("Permissions.SubscriberGroupID").First(subscriberGroupDetail, "id = ? ", subscriberGroupID).Error
	if err != nil {
		errorMessage := fmt.Sprintf("failed to load details of given group id %s, error: %+v", subscriberGroupID, err)
		OSPMLogger.Log.Errorln(errorMessage)
		return models.SubscriberGroup{}, errors.New(errorMessage)
	}

	return subscriberGroupDetail, nil
}

func DeleteSubscriberGroup(subscriberGroupID string) error {
	// we need to find related entities before deletion.
	// Change this part if you know better, god bless ya!
	var permissionSet models.PermissionSets

	err := cockroachdb.DB.Model(&models.PermissionSets{}).
		First(&permissionSet, "subscriber_group_id = ? ", subscriberGroupID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to load details of given group id %s, error: %+v",
			subscriberGroupID, err)
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	deletetionTX := cockroachdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			errorMessage := fmt.Sprintf(
				"failed to delete the given group id %s, error: %+v",
				subscriberGroupID, r)
			OSPMLogger.Log.Errorln(errorMessage)
		}
	}()

	err = deletetionTX.Unscoped().Delete(&models.SubscriberGroup{}).Where("id = ?", subscriberGroupID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete group step, error: %+v",
			subscriberGroupID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	err = deletetionTX.Unscoped().Delete(&models.PermissionSets{}).Where("id = ?", permissionSet.ID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete permission set step, error: %+v",
			subscriberGroupID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	err = deletetionTX.Unscoped().Delete(&models.OrganizationalLevelPerms{}).Where("permission_set_id = ?", permissionSet.ID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete organizationl level perms step, error: %+v",
			subscriberGroupID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	err = deletetionTX.Unscoped().Delete(&models.AccessLevelPerms{}).Where("permission_set_id = ?", permissionSet.ID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete access level perms step, error: %+v",
			subscriberGroupID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	err = deletetionTX.Unscoped().Delete(&models.SubscriberLevelPerms{}).Where("permission_set_id = ?", permissionSet.ID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete subscriber level perms step, error: %+v",
			subscriberGroupID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	err = deletetionTX.Unscoped().Delete(&models.PaymentLevelPerms{}).Where("permission_set_id = ?", permissionSet.ID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete payment level perms step, error: %+v",
			subscriberGroupID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	err = deletetionTX.Unscoped().Delete(&models.ReportLevelPerms{}).Where("permission_set_id = ?", permissionSet.ID).Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at delete report level perms step, error: %+v",
			subscriberGroupID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	err = deletetionTX.Commit().Error
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to delete the given group id %s at apply step, error: %+v",
			subscriberGroupID, err.Error())
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	OSPMLogger.Log.Infoln("subscriber group id %s successfully deleted", subscriberGroupID)

	return nil
}

func AddNewSubscriberGroup(newSubscriberGroup models.SubscriberGroup) error {
	err := cockroachdb.DB.Create(&newSubscriberGroup)
	if err != nil {
		errorMessage := fmt.Sprintf(
			"failed to add the new subscriber group  %s at apply step, error: %+v",
			newSubscriberGroup.Name, err)
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	OSPMLogger.Log.Infoln("subscriber group %s successfully added. id: %s", newSubscriberGroup.Name, newSubscriberGroup.ID)

	return nil
}

func UpdateSubscriber(newSubscriberGroupDetails models.SubscriberGroup, subscriberGroupID string) error {
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
		OSPMLogger.Log.Errorln(errorMessage)
		return errors.New(errorMessage)
	}

	if oldSubscriberGroupDetail.Name != newSubscriberGroupDetails.Name {
		err = cockroachdb.DB.Model(&models.SubscriberGroup{}).
			Update("subscriber_group_name", newSubscriberGroupDetails.Name).
			Where("id = ?", subscriberGroupID).Error
		if err != nil {
			errorMessage := fmt.Sprintf(
				"failed to update the given group id %s, error: %+v",
				subscriberGroupID, err.Error())
			OSPMLogger.Log.Errorln(errorMessage)
			return errors.New(errorMessage)
		}
	}

	if oldSubscriberGroupDetail.Description != newSubscriberGroupDetails.Description {
		err = cockroachdb.DB.Model(&models.SubscriberGroup{}).
			Update("subscriber_group_description", newSubscriberGroupDetails.Description).
			Where("id = ?", subscriberGroupID).Error
		if err != nil {
			errorMessage := fmt.Sprintf(
				"failed to update the given group id %s, error: %+v",
				subscriberGroupID, err.Error())
			OSPMLogger.Log.Errorln(errorMessage)
			return errors.New(errorMessage)
		}
	}

	if oldSubscriberGroupDetail.Description != newSubscriberGroupDetails.Description {
		err = cockroachdb.DB.Model(&models.SubscriberGroup{}).
			Update("subscriber_group_description", newSubscriberGroupDetails.Description).
			Where("id = ?", subscriberGroupID).Error
		if err != nil {
			errorMessage := fmt.Sprintf(
				"failed to update the given group id %s, error: %+v",
				subscriberGroupID, err.Error())
			OSPMLogger.Log.Errorln(errorMessage)
			return errors.New(errorMessage)
		}
	}

	OSPMLogger.Log.Infoln("subscriber group %s successfully updated. id: %s", newSubscriberGroupDetails.Name, subscriberGroupID)

	return nil

}
