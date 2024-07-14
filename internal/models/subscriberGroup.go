package models

import "gorm.io/gorm"

type SubscriberGroup struct {
	gorm.Model
	ID             string         `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"subscriber_group_id"`
	Name           string         `gorm:"not null;index;unique" json:"subscriber_group_name"`
	Description    string         `gorm:"" json:"subscriber_group_description"`
	Permissions    PermissionSets `gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subscriber_group_permissions"`
	OrganizationID string         `gorm:"type:uuid;not null;index;" json:"organization_id"`
}

// SubscriberGroupMinimal is use in API calls which contains a minimal version of
// subscriber group that can be used in listing or addressing
type SubscriberGroupMinimal struct {
	ID   string `json:"subscriber_group_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	Name string `json:"subscriber_group_name" example:"sample group"`
}
