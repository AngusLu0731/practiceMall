package model

import (
	_ "gorm.io/gorm"
)

type Auth struct {
	Id          int
	ModuleName  string //模組名
	ActionName  string //動作名
	Type        int    //捷點類型 :  1、表示模組    2、表示目錄   3、操作
	Url         string //router跳轉地址
	ModuleId    int    //此module_id和當前模型的_id關聯      module_id= 0 表示模組
	Sort        int
	Description string
	Status      int
	AddTime     int
	AuthItem    []Auth `gorm:"foreignkey:ModuleId;association_foreignkey:Id"`
	Checked     bool   `gorm:"-"` // 忽略本字段
}

func (Auth) TableName() string {
	return "auth"
}
