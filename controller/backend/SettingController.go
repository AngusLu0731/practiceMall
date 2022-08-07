package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
)

type SettingController struct {
	BaseController
}

func NewSettingController() *SettingController {
	return &SettingController{}
}

func (c *SettingController) Get(ctx *gin.Context) {
	settingIntoHTML := make(map[string]interface{})
	settingIntoHTML = IntoHTML
	setting := model.Setting{}
	model.DB.First(&setting)
	settingIntoHTML["setting"] = setting
	ctx.HTML(http.StatusOK, "setting_index.html", settingIntoHTML)
}

func (c *SettingController) GoEdit(ctx *gin.Context) {
	setting := model.Setting{}
	model.DB.First(&setting)
	siteLogo, err := c.LocalUploadImage(ctx, ("site_logo"))
	if len(siteLogo) > 0 && err == nil {
		setting.SiteLogo = siteLogo
	}
	noPicture, err := c.LocalUploadImage(ctx, "no_picture")
	if len(noPicture) > 0 && err == nil {
		setting.NoPicture = noPicture
	}
	err = model.DB.Where("id=1").Save(&setting).Error
	if err != nil {
		c.Error(ctx, "修改數據失敗", "/setting")
		return
	}
	c.Success(ctx, "修改數據成功", "/setting")
}
