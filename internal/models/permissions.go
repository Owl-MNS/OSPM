package models

import "gorm.io/gorm"

type PermissionSets struct {
	gorm.Model
	ID                     string                   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"subscriber_group_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	OrganizationLevelPerms OrganizationalLevelPerms `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organizational_level_perms"`
	AccessLevelPerms       AccessLevelPerms         `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"access_level_perms"`
	SubscriberLevelPerms   SubscriberLevelPerms     `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subscriber_level_perms"`
	PaymentLevelPerms      PaymentLevelPerms        `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment_level_perms"`
	ReportLevelPerms       ReportLevelPerms         `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"report_level_perms"`
	OrganizationID         string                   `gorm:"type:uuid;not null;index;" json:"organization_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
}

type OrganizationalLevelPerms struct {
	gorm.Model
	ID              string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"organizational_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `gorm:"type:uuid;not null;index;" json:"permission_set_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
}

type AccessLevelPerms struct {
	gorm.Model
	ID              string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"access_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `gorm:"type:uuid;not null;index;" json:"permission_set_id"`
}

type SubscriberLevelPerms struct {
	gorm.Model
	ID              string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"subscriber_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `gorm:"type:uuid;not null;index;" json:"permission_set_id"`
}

type PaymentLevelPerms struct {
	gorm.Model
	ID              string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"payment_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `gorm:"type:uuid;not null;index;" json:"permission_set_id"`
}

type ReportLevelPerms struct {
	gorm.Model
	ID              string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"report_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `gorm:"type:uuid;not null;index;" json:"permission_set_id"`
}
