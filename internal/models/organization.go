package models

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	ID                       string              `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"organization_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	Details                  OrganizationDetails `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization_details"`
	Owner                    OrganizationOwner   `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization_owner"`
	Balance                  float64             `gorm:"not null;index" json:"balance"`
	AllowNagativeBalance     bool                `gorm:"not null;index" json:"allow_negative_balance"`
	NegativeBalanceThreshold float64             `gorm:"not null;index" json:"negative_balance_threshold"`
}

type OrganizationDetails struct {
	gorm.Model
	ID             string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"organization_details_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	Name           string `gorm:"index;unique" json:"name"`
	Address        string `gorm:"index" json:"address"`
	Email          string `gorm:"index;unique" json:"email"`
	Mobile         string `gorm:"index;unique" json:"mobile"`
	Phone          string `gorm:"index" json:"phone"`
	OrganizationID string `gorm:"type:uuid;not null;index" json:"organization_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
}

type OrganizationOwner struct {
	gorm.Model
	ID              string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"organization_owner_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	Type            string `gorm:"index"` //valid values are: legal, individual
	Name            string `gorm:"index;unique" json:"name"`
	Address         string `gorm:"index" json:"address"`
	Email           string `gorm:"not null;index;unique" json:"email"`
	Mobile          string `gorm:"index;unique" json:"mobile"`
	Phone           string `gorm:"index" json:"phone"`
	LegalNationalID string `gorm:"index" json:"legal_national_id"`
	OrganizationID  string `gorm:"type:uuid;not null;index" json:"organization_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
}

// This model is used while listing the organizations
type OrganizationShortInfo struct {
	ID   string `json:"organization_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"` // This field determines the unique id of the organization. The id is in uuid v4 format
	Name string `json:"organization_name" example:"sample organization"`                // This field determines the name of the organization
}

// The following models are used to represent the raw details of the organization
// in API responses to avoid expose unnecessary details
// Start
type OrganizationResponse struct {
	ID                       string                      `json:"organization_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	Details                  OrganizationDetailsResponse `json:"organization_details"`
	Owner                    OrganizationOwnerResponse   `json:"organization_owner"`
	Balance                  float64                     `json:"balance"`
	AllowNagativeBalance     bool                        `json:"allow_negative_balance"`
	NegativeBalanceThreshold float64                     `json:"negative_balance_threshold"`
}

type OrganizationDetailsResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Mobile  string `json:"mobile"`
	Phone   string `json:"phone"`
}

type OrganizationOwnerResponse struct {
	Type            string `json:"type"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	Email           string `json:"email"`
	Mobile          string `json:"mobile"`
	Phone           string `json:"phone"`
	LegalNationalID string `json:"legal_national_id"`
}

// End of Reponse models!
