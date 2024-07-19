package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	ID                 string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	SubscriberGroupID  string `gorm:"type:uuid;not null;index;" json:"subscriber_group_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7b"`
	PermissionName     string `gorm:"index;" json:"permission_name" example:"CAN_VIEW_PAYMENT_HISTORY"`
	PermissionValue    string `gorm:"" json:"permission_value" example:"yes"`
	PermissionCategory string `gorm:"index;" json:"permission_category" example:"REPORT_LEVEL"`
}

// ##########################
// #	Swagger/API Models	#
// ##########################
// The following models are used for swagger documentation
type PermissionAPI struct {
	PermissionName  string `json:"permission" example:"CAN_VIEW_PAYMENT_HISTORY"`
	PermissionValue string `json:"value" example:"yes"`
}
