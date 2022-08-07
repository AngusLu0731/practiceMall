package model

import (
	_ "gorm.io/gorm"
)

type ProductAttr struct {
	Id              int
	ProductId       int
	AttributeCateId int
	AttributeId     int
	AttributeTitle  string
	AttributeType   int
	AttributeValue  string
	Sort            int
	AddTime         int
	Status          int
}

func (ProductAttr) TableName() string {
	return "product_attr"
}
