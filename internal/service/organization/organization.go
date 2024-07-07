package organization

import (
	"errors"
	"fmt"
	"ospm/internal/models"
	"ospm/internal/repository/database/cockroachdb"
	OSPMLogger "ospm/internal/service/log"
)

// GetOrganizationList returns a list of organizations in shortened format
func GetOrganizationList() ([]models.OrganizationShortInfo, error) {
	organizationList := []models.Organization{}
	result := cockroachdb.DB.Find(&organizationList)
	if result.Error != nil {
		errorMessage := fmt.Sprintf("failed to get list of organization, error: %s", result.Error)
		OSPMLogger.Log.Errorln(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return Shorten(organizationList), nil

}

// GetOrganizationDetails gets the name of the desired organization name and returns the
// details for the given name. Note that the accress credentials are hidden and to check the credentials
// another endpoint should be called
func GetOrganizationDetails(organizationName string, organizationID string) (models.Organization, error) {
	var organization models.Organization

	query := cockroachdb.DB.Preload("Details").Preload("Owner").
		Joins("left join organization_details on organization_details.organization_id = organizations.id")

	if organizationID != "" {
		query = query.Where("organizations.id = ?", organizationID)
	}
	if organizationName != "" {
		query = query.Where("organization_details.name = ?", organizationName)
	}

	err := query.First(&organization).Error
	if err != nil {
		return models.Organization{}, err
	}

	return organization, nil

}

// Shorten gets a list of organizations and returns a list of organizations just including
// ID and Name
func Shorten(organizations []models.Organization) []models.OrganizationShortInfo {
	shortList := []models.OrganizationShortInfo{}
	for _, organization := range organizations {
		shortInfo := models.OrganizationShortInfo{}
		shortInfo.ID = organization.ID
		shortInfo.Name = organization.Details.Name
		shortList = append(shortList, shortInfo)
	}

	return shortList
}
