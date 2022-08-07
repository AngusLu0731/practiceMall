package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"practiceMall/config"
	"practiceMall/model"
	"practiceMall/util"
	"strconv"
	"strings"
)

type BaseController struct {
}

var IntoHTML map[string]interface{}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (c *BaseController) Success(ctx *gin.Context, message string, redirect string) {
	IntoHTML = make(map[string]interface{})
	IntoHTML["Message"] = message
	if strings.Contains(redirect, "http") {
		IntoHTML["Redirect"] = redirect
	} else {
		IntoHTML["Redirect"] = "/" + config.Conf.AdminPath + redirect
	}
	ctx.HTML(http.StatusOK, "backend_success.html", IntoHTML)
}

func (c *BaseController) Error(ctx *gin.Context, message string, redirect string) {
	IntoHTML = make(map[string]interface{})
	IntoHTML["Message"] = message
	if strings.Contains(redirect, "http") {
		IntoHTML["Redirect"] = redirect
	} else {
		IntoHTML["Redirect"] = "/" + config.Conf.AdminPath + redirect
	}
	ctx.HTML(http.StatusOK, "backend_error.html", IntoHTML)
}
func (c *BaseController) Goto(ctx *gin.Context, redirect string) {
	ctx.Redirect(http.StatusFound, "/"+config.Conf.AdminPath+redirect)
}

func (c *BaseController) LocalUploadImage(ctx *gin.Context, picName string) (string, error) {
	//1.獲取文件
	f, h, err := ctx.Request.FormFile(picName)
	if err != nil {
		return "", err
	}
	//2.defer 關閉文件流
	defer f.Close()

	//3.獲取副檔名，判斷文件是否正確
	extName := path.Ext(h.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}
	if _, ok := allowExtMap[extName]; !ok {
		return "", err
	}

	//4.創建檔案保存目錄
	day := util.FormatDay()
	dir := "static/upload/" + day

	if err := os.MkdirAll(dir, 0666); err != nil {
		return "", err
	}

	//5.生成檔案名稱
	fileUnixName := strconv.FormatInt(util.GetUnixNano(), 10)
	saveDir := path.Join(dir, fileUnixName+extName)

	//6.保存檔案
	ctx.SaveUploadedFile(h, saveDir)
	return saveDir, nil
}

func (c *BaseController) GetSetting() model.Setting {
	setting := model.Setting{Id: 1}
	model.DB.First(&setting)
	return setting
}
