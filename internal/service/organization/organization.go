package organization

import (
	"errors"
	"fmt"
	"ospm/config"
	"ospm/internal/models"
	"ospm/internal/repository/database/cockroachdb"
	"ospm/internal/service/complementary"
	"ospm/internal/service/logger"
	"strings"
)

// List returns a list of organizations in shortened format.
// Organizations that are hard deleted will not be listed
func List() ([]models.OrganizationShortInfo, error) {
	organizationList := []models.Organization{}
	result := cockroachdb.DB.Preload("Details").Find(&organizationList)
	if result.Error != nil {
		errorMessage := fmt.Sprintf("failed to get list of organization, error: %s", result.Error)
		logger.OSPMLogger.Errorln(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return Shorten(organizationList), nil
}

// ListAll returns a list of organizations in shortened format.
// Organizations that are hard deleted will be listed
func ListAll() ([]models.OrganizationShortInfo, error) {
	organizationList := []models.Organization{}
	result := cockroachdb.DB.Unscoped().Preload("Details").Find(&organizationList)
	if result.Error != nil {
		errorMessage := fmt.Sprintf("failed to get list of organization, error: %s", result.Error)
		logger.OSPMLogger.Errorln(errorMessage)
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
		logger.OSPMLogger.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	if err := cockroachdb.DB.Create(&newOrganization).Error; err != nil {
		errorMessage := fmt.Sprintf("the new organization can not be created, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	return newOrganization.ID, nil
}

// SoftDelete deletes the desired organization and does not impact
// the other related entities like subscriber, permissions and etc.
// The delete action happens in soft mode
func SoftDelete(organizationID string, organizationName string) error {
	var organization models.Organization

	query := cockroachdb.DB.Joins("left join organization_details on organization_details.organization_id = organizations.id")

	if organizationID != "" {
		query = query.Where("organizations.id = ?", organizationID)
	} else if organizationName != "" {
		query = query.Where("organization_details.name = ?", organizationName)
	}

	err := query.First(&organization).Error
	if err != nil {
		errorMessage := fmt.Sprintf("failed to find organization to delete, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	// Start a transaction
	tx := cockroachdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Perform the delete operation with cascading deletes with HARD DELETE Enabled!
	if err := tx.Select("Details", "Owner").Delete(&organization).Error; err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("failed to delete organization and related records, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("failed to commit transaction, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	return nil
}

// HardDelete deletes the desired organization and does not impact
// the other related entities like subscriber, permissions and etc.
// The delete action happens in hard mode
func HardDelete(organizationID string, organizationName string) error {
	var organization models.Organization

	query := cockroachdb.DB.Unscoped().Joins("left join organization_details on organization_details.organization_id = organizations.id")

	if organizationID != "" {
		query = query.Where("organizations.id = ?", organizationID)
	} else if organizationName != "" {
		query = query.Where("organization_details.name = ?", organizationName)
	}

	err := query.First(&organization).Error
	if err != nil {
		errorMessage := fmt.Sprintf("failed to find organization to delete, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	// Start a transaction
	tx := cockroachdb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Perform the delete operation with cascading deletes with HARD DELETE Enabled!
	if err := tx.Unscoped().Select("Details", "Owner").Delete(&organization).Error; err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("failed to delete organization and related records, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("failed to commit transaction, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	return nil
}

// Recover truncates the deleted_at field from the database which
// recovers the organization from soft delete
func Recover(organizationID string, organizationName string) error {
	var organization models.Organization

	query := cockroachdb.DB.Unscoped().Joins("left join organization_details on organization_details.organization_id = organizations.id")

	if organizationID != "" {
		query = query.Where("organizations.id = ?", organizationID)
	} else if organizationName != "" {
		query = query.Where("organization_details.name = ?", organizationName)
	}

	err := query.First(&organization).Error
	if err != nil {
		errorMessage := fmt.Sprintf("failed to find organization to delete, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	// Start the transaction
	tx := cockroachdb.DB.Begin()

	// Restore the organization
	if err := tx.Unscoped().Model(&models.Organization{}).Where("id = ?", organization.ID).Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("failed to recover organization from soft delete, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	// Restore organization details
	if err := tx.Unscoped().Model(&models.OrganizationDetails{}).Where("organization_id = ?", organization.ID).Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("failed to recover organization from soft delete, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	// Restore related owner
	if err := tx.Unscoped().Model(&models.OrganizationOwner{}).Where("organization_id = ?", organization.ID).Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("failed to recover organization from soft delete, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		errorMessage := fmt.Sprintf("failed to recover organization from soft delete, error: %+v", err)
		logger.OSPMLogger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	return nil
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
func DetailsCheck(organizationDetails *models.Organization) error {
	var err error

	if organizationDetails.Balance != 0 {
		errorMessage := fmt.Sprintf(
			"organization balance can not accept any values but 0 while creating the organization. given value is: %f",
			organizationDetails.Balance)
		err = errors.New(errorMessage)
	}

	if organizationDetails.AllowNagativeBalance {
		errorMessage := fmt.Sprintf(
			"organization AllowNagativeBalance can not be true while creating the organization. given value is: %v",
			organizationDetails.AllowNagativeBalance)
		err = errors.New(errorMessage)
	}

	if organizationDetails.NegativeBalanceThreshold != 0 {
		errorMessage := fmt.Sprintf(
			"organization NegativeBalanceThreshold can not accept any values but 0 while creating the organization. given value is: %f",
			organizationDetails.NegativeBalanceThreshold)
		err = errors.New(errorMessage)
	}

	if organizationDetails.Details.Name == "" {
		errorMessage := fmt.Sprintf(
			"organization Name can not be empty while creating the organization. given value is: %s",
			organizationDetails.Details.Name)
		err = errors.New(errorMessage)
	}

	if organizationDetails.Owner.Email == "" {
		errorMessage := fmt.Sprintf(
			"organization's Owner email address can not be empty while creating the organization. given value is: %s",
			organizationDetails.Owner.Email)
		err = errors.New(errorMessage)
	}

	if organizationDetails.Owner.Mobile == "" {
		errorMessage := fmt.Sprintf(
			"organization's Owner Mobile can not be empty while creating the organization. given value is: %s",
			organizationDetails.Owner.Mobile)
		err = errors.New(errorMessage)
	}

	if !(organizationDetails.Owner.Type == "legal" || organizationDetails.Owner.Type == "individual") {
		errorMessage := fmt.Sprintf(
			"organization's Owner typ should be either individual or legal while creating the organization. given value is: %s",
			organizationDetails.Owner.Type)
		err = errors.New(errorMessage)
	}

	if organizationDetails.Owner.LegalNationalID == "" {
		errorMessage := fmt.Sprintf(
			"organization's Owner Legal National ID can not be empty while creating the organization. given value is: %s",
			organizationDetails.Owner.LegalNationalID)
		err = errors.New(errorMessage)
	}

	if err != nil {
		errorMessage := fmt.Sprintf("new organization details are wrong. error: %+v", err)
		return errors.New(errorMessage)
	}
	return nil
}

// Clean can be used to remove database related items from the results returned from the
// DB query like created_at, deleted_at and etc.
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

// ClientIPCanSoftDeleteOrganization gets the client's IP and checks it among
// the permited IPs. If the client's ip is whitelisted, returns true
func ClientIPCanSoftDeleteOrganization(clientIP string) bool {

	// Check if the client's IP is in the allowed list or ranges
	for _, allowedIP := range strings.Split(config.OSPM.ClientPolicies.OrganizationSoftDeleteWhiteListedIPs, ",") {
		if complementary.IPRangeCotains(clientIP, allowedIP) {
			return true
		}
	}

	return false
}

// ClientIPCanSoftDeleteOrganization gets the client's IP and checks it among
// the permited IPs. If the client's ip is whitelisted, returns true
func ClientIPCanHardDeleteOrganization(clientIP string) bool {

	// Check if the client's IP is in the allowed list or ranges
	for _, allowedIP := range strings.Split(config.OSPM.ClientPolicies.OrganizationHardDeleteWhiteListedIPs, ",") {
		if complementary.IPRangeCotains(clientIP, allowedIP) {
			return true
		}
	}

	return false
}

// ClientIPCanListAllOrganization gets the client's IP and checks it among
// the permited IPs. If the client's ip is whitelisted, returns true
// whitelisted ips can list all of the organizations including soft deleteds
func ClientIPCanListAllOrganization(clientIP string) bool {

	// Check if the client's IP is in the allowed list or ranges
	for _, allowedIP := range strings.Split(config.OSPM.ClientPolicies.ListAllOrganizationWhiteListedIPs, ",") {
		if complementary.IPRangeCotains(clientIP, allowedIP) {
			return true
		}
	}

	return false
}

// ClientIPCanUndoOrganizationSoftDelete gets the client ip and checks it
// among permitted IPs. If the client's ip is whitelisted, returns true
func ClientIPCanUndoOrganizationSoftDelete(clientIP string) bool {

	// Check if the client's IP is in the allowed list or ranges
	for _, allowedIP := range strings.Split(config.OSPM.ClientPolicies.UndoOrganizationSoftDeleteWhiteListedIPs, ",") {
		if complementary.IPRangeCotains(clientIP, allowedIP) {
			return true
		}
	}

	return false
}
