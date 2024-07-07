package models

import "gorm.io/gorm"

type Subscriber struct {
	gorm.Model
	ID                string            `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"subscriber_id"`
	Details           SubscriberDetails `gorm:"foreignKey:SubscriberID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subscriber_details"`
	Credentials       Credentials       `gorm:"foreignKey:SubscriberID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subscriber_credentials"`
	OrganizationID    string            `gorm:"type:uuid;not null;index;" json:"organization_id"`
	SubscriberGroupID string            `gorm:"type:uuid;not null;index" json:"subsdriber_group_id"`
	// Offers            []Offer
}

type SubscriberDetails struct {
	gorm.Model
	ID           string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"subscriber_details_id"`
	Name         string `gorm:"not null;index;unique" json:"subscriber_name"`
	Email        string `gorm:"not null;index;unique" json:"subscriber_email"`
	NationalID   string `gorm:"index;unique" json:"subscriber_national_id"`
	PassportID   string `gorm:"index;unique" json:"passport_id"`
	Mobile       string `gorm:"not null;index;unique" json:"subscriber_mobile"`
	Phone        string `gorm:"index;unique" json:"subscriber_phone"`
	SubscriberID string `gorm:"type:uuid;not null;index" json:"subsdriber_id"`
}

type Credentials struct {
	gorm.Model
	ID                  string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"subscriber_credential_id"`
	Username            string `gorm:"not null;index;unique" json:"subscriber_username"`
	Password            string `gorm:"not null;index;unique" json:"subscriber_password"`
	AuthenticationToken string `gorm:"not null;index;unique" json:"subscriber_authentication_token"`
	SubscriberID        string `gorm:"type:uuid;not null;index" json:"subsdriber_id"`
}
