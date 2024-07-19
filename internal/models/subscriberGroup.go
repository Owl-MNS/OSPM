package models

import "gorm.io/gorm"

type SubscriberGroup struct {
	gorm.Model
	ID             string       `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name           string       `gorm:"not null;index;uniqueIndex:org_name_idx;" json:"subscriber_group_name"`
	Description    string       `gorm:"" json:"subscriber_group_description"`
	Permissions    []Permission `gorm:"foreignKey:SubscriberGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subscriber_group_permissions"`
	OrganizationID string       `gorm:"type:uuid;not null;uniqueIndex:org_name_idx;" json:"organization_id"`
}

func (sg *SubscriberGroup) Beautify() SubscriberGroupAPI {
	beautified := SubscriberGroupAPI{
		ID:             sg.ID,
		Name:           sg.Name,
		Description:    sg.Description,
		OrganizationID: sg.OrganizationID,
		Permissions:    map[string][]PermissionAPI{},
	}

	for _, perm := range sg.Permissions {
		beautified.Permissions[perm.PermissionCategory] = append(beautified.Permissions[perm.PermissionCategory], PermissionAPI{
			PermissionName:  perm.PermissionName,
			PermissionValue: perm.PermissionValue,
		})
	}

	return beautified
}

// Absorb gets an object of SubscriberGroupAPI and loads its content into the current object's fields!
func (sg *SubscriberGroup) Absorb(apiVersion SubscriberGroupAPI) {
	sg.Name = apiVersion.Name
	sg.Description = apiVersion.Description
	sg.OrganizationID = apiVersion.OrganizationID

	for permCategory, permSettings := range apiVersion.Permissions {
		for _, perm := range permSettings {
			newPerm := Permission{}
			newPerm.PermissionCategory = permCategory
			newPerm.PermissionName = perm.PermissionName
			newPerm.PermissionValue = perm.PermissionValue
			sg.Permissions = append(sg.Permissions, newPerm)
		}
	}
}

// ##########################
// #	Swagger/API Models	#
// ##########################
// The following models are used for swagger documentation
type SubscriberGroupAPI struct {
	ID             string                     `json:"subscriber_group_id"`
	Name           string                     `json:"subscriber_group_name"`
	Description    string                     `json:"subscriber_group_description"`
	Permissions    map[string][]PermissionAPI `json:"subscriber_group_permissions"` // This field gets a list of permission categorized per permission category
	OrganizationID string                     `json:"organization_id"`
}

type SubscriberGroupCreateResponse struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Id      string `json:"id"`
}

// SubscriberGroupMinimal is use in API calls which contains a minimal version of
// subscriber group that can be used in listing or addressing
type SubscriberGroupMinimal struct {
	ID   string `json:"subscriber_group_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	Name string `json:"subscriber_group_name" example:"sample group"`
}

type AddSubscriberGroupAPI struct {
	PermissionName string `json:""`
}
