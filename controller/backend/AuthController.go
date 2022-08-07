package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
	"strconv"
)

type AuthController struct {
	BaseController
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) Get(ctx *gin.Context) {
	AuthIntoHTML := make(map[string]interface{})
	AuthIntoHTML = IntoHTML
	auth := []model.Auth{}
	model.DB.Preload("AuthItem").Where("module_id = 0").Find(&auth)
	AuthIntoHTML["authList"] = auth
	ctx.HTML(http.StatusOK, "auth_index.html", AuthIntoHTML)
}

func (c *AuthController) Add(ctx *gin.Context) {
	AuthIntoHTML := make(map[string]interface{})
	AuthIntoHTML = IntoHTML
	auth := []model.Auth{}
	model.DB.Where("module_id=0").Find(&auth)
	AuthIntoHTML["authList"] = auth
	ctx.HTML(http.StatusOK, "auth_add.html", AuthIntoHTML)
}

func (c *AuthController) GoAdd(ctx *gin.Context) {
	moduleName := ctx.PostForm("module_name")
	actionName := ctx.PostForm("action_name")
	url := ctx.PostForm("url")
	description := ctx.PostForm("description")
	iType, err1 := strconv.Atoi(ctx.PostForm("type"))
	moduleId, err2 := strconv.Atoi(ctx.PostForm("module_id"))
	sort, err3 := strconv.Atoi(ctx.PostForm("sort"))
	status, err4 := strconv.Atoi(ctx.PostForm("status"))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Error(ctx, "參數錯誤", "/auth/add")
		return
	}
	auth := model.Auth{
		ModuleName:  moduleName,
		ActionName:  actionName,
		Type:        iType,
		Url:         url,
		ModuleId:    moduleId,
		Sort:        sort,
		Description: description,
		Status:      status,
	}
	err := model.DB.Create(&auth).Error
	if err != nil {
		c.Error(ctx, "添加失敗", "auth/add")
		return
	}
	c.Success(ctx, "添加成功", "/auth")
}

func (c *AuthController) Edit(ctx *gin.Context) {
	AuthIntoHTML := make(map[string]interface{})
	AuthIntoHTML = IntoHTML
	id, err1 := strconv.Atoi(ctx.Query("id"))
	if err1 != nil {
		c.Error(ctx, "參數錯誤", "/auth")
		return
	}
	auth := model.Auth{Id: id}
	model.DB.Find(&auth)
	AuthIntoHTML["auth"] = auth
	authList := []model.Auth{}
	model.DB.Where("module_id = 0").Find(&authList)
	AuthIntoHTML["authList"] = authList
	ctx.HTML(http.StatusOK, "auth_edit.html", AuthIntoHTML)
}

func (c *AuthController) GoEdit(ctx *gin.Context) {
	moduleName := ctx.PostForm("module_name")
	actionName := ctx.PostForm("action_name")
	url := ctx.PostForm("url")
	description := ctx.PostForm("description")
	iType, err1 := strconv.Atoi(ctx.PostForm("type"))
	moduleId, err2 := strconv.Atoi(ctx.PostForm("module_id"))
	sort, err3 := strconv.Atoi(ctx.PostForm("sort"))
	status, err4 := strconv.Atoi(ctx.PostForm("status"))
	id, err5 := strconv.Atoi(ctx.PostForm("id"))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		c.Error(ctx, "參數錯誤", "/auth")
		return
	}

	auth := model.Auth{Id: id}
	model.DB.Find(&auth)
	auth.ModuleName = moduleName
	auth.Type = iType
	auth.ActionName = actionName
	auth.Url = url
	auth.ModuleId = moduleId
	auth.Sort = sort
	auth.Description = description
	auth.Status = status
	err := model.DB.Save(&auth).Error
	if err != nil {
		c.Error(ctx, "修改權限失敗", "/auth/edit?id="+strconv.Itoa(id))
		return
	}
	c.Success(ctx, "修改權限成功", "/auth")
}

func (c *AuthController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/auth")
		return
	}
	auth := model.Auth{Id: id}
	model.DB.Find(&auth)
	if auth.ModuleId == 0 {
		auth2 := []model.Auth{}
		model.DB.Where("module_id=?", auth.Id).Find(&auth2)
		if len(auth2) > 0 {
			c.Error(ctx, "請刪除當前群組中底下的操作再進行刪除", "/auth")
			return
		}
	}
	model.DB.Delete(&auth)
	c.Success(ctx, "刪除成功", "/auth")
}
