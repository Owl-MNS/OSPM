package models

import "gorm.io/gorm"

type ProductOfferingSpecification struct {
	gorm.Model
	ID          string `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"product_offering_specification_id"`
	Name        string `gorm:"not null;index;unique" json:"product_offering_specification_name"`
	PersianName string `gorm:"index;unique" json:"product_offering_specification_persian_name"`
	Type        string `gorm:"not null;" json:"product_offering_specification_type"` // valid values: product, service
	Description string `json:"product_offering_specification_description"`
}
