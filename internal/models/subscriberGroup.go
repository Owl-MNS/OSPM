package models

import "gorm.io/gorm"

type SubscriberGroup struct {
	gorm.Model
	ID             string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"subscriber_group_id"`
	Name           string `gorm:"not null;index;unique" json:"subscriber_group_name"`
	Description    string `gorm:"" json:"subscriber_group_description"`
	Permissions    string `gorm:"type:uuid;" json:"subscriber_group_permissions_id"`
	OrganizationID string `gorm:"type:uuid;not null;index;" json:"organization_id"`
}
