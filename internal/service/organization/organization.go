package organization

import (
	"errors"
	"fmt"
	"ospm/internal/models"
	"ospm/internal/repository/database/cockroachdb"
	OSPMLogger "ospm/internal/service/log"
)

// List returns a list of organizations in shortened format
func List() ([]models.OrganizationShortInfo, error) {
	organizationList := []models.Organization{}
	result := cockroachdb.DB.Preload("Details").Find(&organizationList)
	if result.Error != nil {
		errorMessage := fmt.Sprintf("failed to get list of organization, error: %s", result.Error)
		OSPMLogger.Log.Errorln(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return Shorten(organizationList), nil

}

// Details gets the name of the desired organization name and returns the
// details for the given name. Note that the accress credentials are hidden and to check the credentials
// another endpoint should be called
func Details(organizationName string, organizationID string) (models.Organization, error) {
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

// New gets the new organization details and adds it into the database then returns
// the new added organization's ID. in case of any issue while adding the new organization
// it returns an error
func New(newOrganization models.Organization) (newOrganzationID string, err error) {
	if err := DetailsCheck(&newOrganization); err != nil {
		errorMessage := fmt.Sprintf("the new organization can not be created, error: %+v", err)
		OSPMLogger.Log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	if err := cockroachdb.DB.Create(&newOrganization).Error; err != nil {
		errorMessage := fmt.Sprintf("the new organization can not be created, error: %+v", err)
		OSPMLogger.Log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	return newOrganization.ID, nil
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

// DetailsCheck checks the given new organization details an validates the given values
// Since the given values have met the creation policies, the new organization can be created
// by returning nil as error otherwise the error determines what is wrong with the new given information
func DetailsCheck(newOrganization *models.Organization) error {
	var err error

	if newOrganization.Balance != 0 {
		errorMessage := fmt.Sprintf("organization balance can not accept any values but 0 while creating the organization. given value is: %f", newOrganization.Balance)
		err = errors.New(errorMessage)
	}

	if newOrganization.AllowNagativeBalance {
		errorMessage := fmt.Sprintf("organization AllowNagativeBalance can not be true while creating the organization. given value is: %v", newOrganization.AllowNagativeBalance)
		err = errors.New(errorMessage)
	}

	if newOrganization.NegativeBalanceThreshold != 0 {
		errorMessage := fmt.Sprintf("organization NegativeBalanceThreshold can not accept any values but 0 while creating the organization. given value is: %f", newOrganization.NegativeBalanceThreshold)
		err = errors.New(errorMessage)
	}

	if newOrganization.Details.Name == "" {
		errorMessage := fmt.Sprintf("organization Name can not be empty while creating the organization. given value is: %s", newOrganization.Details.Name)
		err = errors.New(errorMessage)
	}

	if newOrganization.Owner.Email == "" {
		errorMessage := fmt.Sprintf("organization's Owner email address can not be empty while creating the organization. given value is: %s", newOrganization.Owner.Email)
		err = errors.New(errorMessage)
	}

	if newOrganization.Owner.Mobile == "" {
		errorMessage := fmt.Sprintf("organization's Owner Mobile can not be empty while creating the organization. given value is: %s", newOrganization.Owner.Mobile)
		err = errors.New(errorMessage)
	}

	if !(newOrganization.Owner.Type == "legal" || newOrganization.Owner.Type == "individual") {
		errorMessage := fmt.Sprintf("organization's Owner typ should be either individual or legal while creating the organization. given value is: %s", newOrganization.Owner.Type)
		err = errors.New(errorMessage)
	}

	if newOrganization.Owner.LegalNationalID == "" {
		errorMessage := fmt.Sprintf("organization's Owner Legal National ID can not be empty while creating the organization. given value is: %s", newOrganization.Owner.LegalNationalID)
		err = errors.New(errorMessage)
	}

	if err != nil {
		errorMessage := fmt.Sprintf("new organization details are wrong. error: %+v", err)
		return errors.New(errorMessage)
	}
	return nil
}

func Clean(organization *models.Organization) models.OrganizationResponse {
	return models.OrganizationResponse{
		ID:                       organization.ID,
		Balance:                  organization.Balance,
		AllowNagativeBalance:     organization.AllowNagativeBalance,
		NegativeBalanceThreshold: organization.NegativeBalanceThreshold,
		Details: models.OrganizationDetailsResponse{
			Name:    organization.Details.Name,
			Address: organization.Details.Address,
			Email:   organization.Details.Email,
			Mobile:  organization.Details.Mobile,
			Phone:   organization.Details.Phone,
		},
		Owner: models.OrganizationOwnerResponse{
			Type:            organization.Owner.Type,
			Name:            organization.Owner.Name,
			Address:         organization.Owner.Address,
			Email:           organization.Owner.Email,
			Mobile:          organization.Owner.Mobile,
			Phone:           organization.Owner.Phone,
			LegalNationalID: organization.Owner.LegalNationalID,
		},
	}

}
