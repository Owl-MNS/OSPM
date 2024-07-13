package organization

import (
	"errors"
	"fmt"
	"ospm/internal/models"
	"ospm/internal/repository/database/cockroachdb"
	OSPMLogger "ospm/internal/service/log"
)

// GetOrganizationGroupList get the organization id and returns all groups within the given organiztion
// In the listed group, soft deleted groups are excluded!
func GetOrganizationGroupList(organizationsID string) ([]models.SubscriberGroupMinimal, error) {
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
