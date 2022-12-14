package model

import (
	_ "gorm.io/gorm"
)

type Banner struct {
	Id         int
	Title      string
	BannerType int
	BannerImg  string
	Link       string
	Sort       int
	Status     int
	AddTime    int
}

func (Banner) TableName() string {
	return "banner"
}
