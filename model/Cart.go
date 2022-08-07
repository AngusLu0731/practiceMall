package model

import (
	_ "gorm.io/gorm"
)

type Cart struct {
	Id             int
	Title          string
	Price          float64
	ProductVersion string
	Num            int
	ProductGift    string
	ProductFitting string
	ProductColor   string
	ProductImg     string
	ProductAttr    string
	Checked        bool `gorm:"-"` // 忽略本字段
}

func (Cart) TableName() string {
	return "cart"
}

// CartHasData 判斷購物車內有無當前數據
func CartHasData(cartList []Cart, currentData Cart) bool {
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == currentData.Id &&
			cartList[i].ProductColor == currentData.ProductColor &&
			cartList[i].ProductAttr == currentData.ProductAttr {
			return true
		}
	}
	return false
}
