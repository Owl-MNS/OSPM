package models

import "gorm.io/gorm"

type ProductOffering struct {
	gorm.Model
	ID                      string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"product_offering_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	Name                    string `gorm:"not null;index;unique" json:"product_offering_name"`
	PersianName             string `gorm:"index;unique" json:"product_offering_persian_name"`
	CharacteristicValue     string `gorm:"" json:"product_offering_characteristic_value"`
	CharacteristicValueType string `gorm:"" json:"product_offering_characteristic_value_type"`
	SpecificationID         string `gorm:"type:uuid" json:"product_offering_specification_id" example:"ed83a2ba-c55c-4297-b2ac-df7b02abdd7a"`
	Description             string `json:"product_offering_description"`
}
