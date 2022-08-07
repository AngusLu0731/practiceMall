package model

import (
	_ "gorm.io/gorm"
)

type ProductItemAttr struct {
	Cate string
	List []string
}
