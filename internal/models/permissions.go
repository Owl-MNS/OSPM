package models

import "gorm.io/gorm"

type PermissionSets struct {
	gorm.Model
	ID                     string                   `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"permission_sets_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	OrganizationLevelPerms OrganizationalLevelPerms `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organizational_level_perms"`
	AccessLevelPerms       AccessLevelPerms         `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"access_level_perms"`
	SubscriberLevelPerms   SubscriberLevelPerms     `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subscriber_level_perms"`
	PaymentLevelPerms      PaymentLevelPerms        `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payment_level_perms"`
	ReportLevelPerms       ReportLevelPerms         `gorm:"foreignKey:PermissionSetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"report_level_perms"`
	SubscriberGroupID      string                   `gorm:"type:uuid;not null;index;" json:"subscriber_group_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7b"`
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

// ######################
// #	Swagger Models	#
// ######################
// The following models are used for swagger documentation
type PermissionSetsSwagger struct {
	ID                     string                          `json:"permission_sets_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	OrganizationLevelPerms OrganizationalLevelPermsSwagger `json:"organizational_level_perms"`
	AccessLevelPerms       AccessLevelPermsSwagger         `json:"access_level_perms"`
	SubscriberLevelPerms   SubscriberLevelPermsSwagger     `json:"subscriber_level_perms"`
	PaymentLevelPerms      PaymentLevelPermsSwagger        `json:"payment_level_perms"`
	ReportLevelPerms       ReportLevelPermsSwagger         `json:"report_level_perms"`
	SubscriberGroupID      string                          `json:"subscriber_group_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7b"`
}

// This model is used for swagger documentation
type OrganizationalLevelPermsSwagger struct {
	ID              string `json:"organizational_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `json:"permission_set_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
}

// This model is used for swagger documentation
type SubscriberLevelPermsSwagger struct {
	ID              string `json:"subscriber_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `json:"permission_set_id"`
}

// This model is used for swagger documentation
type AccessLevelPermsSwagger struct {
	ID              string `json:"access_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `json:"permission_set_id"`
}

// This model is used for swagger documentation
type PaymentLevelPermsSwagger struct {
	ID              string `json:"payment_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `json:"permission_set_id"`
}

// This model is used for swagger documentation
type ReportLevelPermsSwagger struct {
	ID              string `json:"report_level_perms_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	PermissionSetID string `json:"permission_set_id"`
}
