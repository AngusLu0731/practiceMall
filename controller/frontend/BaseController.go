package frontend

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/url"
	"practiceMall/config"
	"practiceMall/model"
	"strings"
	"time"
)

type BaseController struct{}

var IntoHTML map[string]interface{}

func (c *BaseController) BaseInit(ctx *gin.Context) {
	IntoHTML = make(map[string]interface{})
	//獲取上方導覽選單
	topMenu := []model.Menu{}
	if hasMenuTop := model.RedisStore.Get("topMenu", &topMenu); hasMenuTop == nil {
		IntoHTML["topMenuList"] = topMenu
	} else {
		model.DB.Where("status=1 AND position=1").Order("sort desc").Find(&topMenu)
		IntoHTML["topMenuList"] = topMenu
		model.RedisStore.Set("topMenu", topMenu, time.Duration(config.Conf.RedisTime))
	}

	//預加載左側分類
	productCate := []model.ProductCate{}
	if hasProductCate := model.RedisStore.Get("productCate", &productCate); hasProductCate == nil {
		IntoHTML["productCateList"] = productCate
	} else {
		model.DB.Preload("ProductCateItem", func(db *gorm.DB) *gorm.DB {
			return db.Where("product_cate.status=1").
				Order("product_cate.sort DESC")
		}).Where("pid=0 AND status=1").Order("sort desc").
			Find(&productCate)
		IntoHTML["productCateList"] = productCate
		model.RedisStore.Set("productCate", productCate, time.Duration(config.Conf.RedisTime))
	}

	//獲取中間Menu的數據
	middleMenu := []model.Menu{}
	if hasMiddleMenu := model.RedisStore.Get("middleMenu", &middleMenu); hasMiddleMenu == nil {
		IntoHTML["middleMenu"] = middleMenu
	} else {
		model.DB.Where("status=1 AND position=2").Order("sort desc").Find(&middleMenu)
		for i := 0; i < len(middleMenu); i++ {
			//獲取關聯商品
			middleMenu[i].Relation = strings.ReplaceAll(middleMenu[i].Relation, "，", ",")
			relation := strings.Split(middleMenu[i].Relation, ",")
			product := []model.Product{}
			model.DB.Where("id in (?)", relation).Limit(6).Order("sort ASC").
				Select("id,title,product_img,price").Find(&product)
			middleMenu[i].ProductItem = product
		}
		IntoHTML["middleMenuList"] = middleMenu
		model.RedisStore.Set("middleMenu", middleMenu, time.Duration(config.Conf.RedisTime))
	}

	//判斷是否登入
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	if len(user.Phone) == 10 {
		str := user.Phone
		IntoHTML["userinfo"] = str
	} else {
		str := "null"
		IntoHTML["userinfo"] = str
	}
	urlPath, _ := url.Parse(ctx.Request.URL.String())
	IntoHTML["pathname"] = urlPath.Path
}
