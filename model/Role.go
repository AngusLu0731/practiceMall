package model

import (
	_ "gorm.io/gorm"
)

type Role struct {
	Id          int
	Title       string
	Description string
	Status      int
	AddTime     int
}

func (Role) TableName() string {
	return "role"
}
