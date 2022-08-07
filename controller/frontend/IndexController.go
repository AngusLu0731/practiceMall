package frontend

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"practiceMall/config"
	"practiceMall/model"
	"time"
)

type IndexController struct {
	BaseController
}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (c *IndexController) Get(ctx *gin.Context) {
	c.BaseInit(ctx)
	indexIntoHTML := make(map[string]interface{})
	indexIntoHTML = IntoHTML

	startTime := time.Now().UnixNano()
	//獲取輪播圖
	banner := []model.Banner{}
	if hasBanner := model.RedisStore.Get("banner", &banner); hasBanner == nil {
		indexIntoHTML["bannerList"] = banner
	} else {
		model.DB.Where("status = 1 AND banner_type = 1").Order("sort desc").Find(&banner)
		indexIntoHTML["bannerList"] = banner
		model.RedisStore.Set("banner", banner, time.Duration(config.Conf.RedisTime))
	}

	//獲取手機商品列表
	phoneProduct := []model.Product{}
	if hasPhone := model.RedisStore.Get("phone", &phoneProduct); hasPhone == nil {
		indexIntoHTML["phoneList"] = phoneProduct
	} else {
		phone := model.GetProductByCategory(1, "hot", 8)
		indexIntoHTML["phoneList"] = phone
		model.RedisStore.Set("phone", phone, time.Duration(config.Conf.RedisTime))
	}

	//獲取電視商品列表
	tvProduct := []model.Product{}
	if hasTv := model.RedisStore.Get("tv", &tvProduct); hasTv == nil {
		indexIntoHTML["tvList"] = tvProduct
	} else {
		tv := model.GetProductByCategory(4, "best", 8)
		indexIntoHTML["tvList"] = tv
		model.RedisStore.Set("tv", tv, time.Duration(config.Conf.RedisTime))
	}
	//結束時間
	endTime := time.Now().UnixNano()
	log.Println("執行時間：", endTime-startTime)
	ctx.HTML(http.StatusOK, "index.html", indexIntoHTML)
}
