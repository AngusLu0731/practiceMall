package model

import (
	_ "gorm.io/gorm"
)

type Order struct {
	Id          int
	OrderId     string
	Uid         int
	AllPrice    float64
	Phone       string
	Name        string
	Address     string
	Zipcode     string
	PayStatus   int //支付狀態： 0 未支付     1 已支付
	PayType     int // 支付類型： 0 LinePay  1 待開發
	OrderStatus int // 訂單狀態： 0 已下單  1 已付款  2 已配貨  3 發貨   4 交易成功  5 退貨  6 取消
	AddTime     int
	OrderItem   []OrderItem `gorm:"foreignkey:OrderId;association_foreignkey:Id"`
}

func (Order) TableName() string {
	return "order"
}
