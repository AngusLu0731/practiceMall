package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
	"practiceMall/util"
	"strconv"
	"strings"
)

type RoleController struct {
	BaseController
}

func NewRoleController() *RoleController {
	return &RoleController{}
}

func (c *RoleController) Get(ctx *gin.Context) {
	RoleIntoHTML := make(map[string]interface{})
	RoleIntoHTML = IntoHTML
	role := []model.Role{}
	model.DB.Find(&role)
	RoleIntoHTML["rolelist"] = role
	ctx.HTML(http.StatusOK, "role_index.html", RoleIntoHTML)
}

func (c *RoleController) Add(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "role_add.html", nil)
}

func (c *RoleController) GoAdd(ctx *gin.Context) {
	title := strings.Trim(ctx.PostForm("title"), "")
	description := strings.Trim(ctx.PostForm("description"), "")
	if title == "" {
		c.Error(ctx, "標題不能為空", "/role/add")
		return
	}
	roleList := []model.Role{}
	model.DB.Where("title=?", title).Find(&roleList)
	if len(roleList) != 0 {
		c.Error(ctx, "該部門已存在", "/role/add")
		return
	}
	role := model.Role{
		Title:       title,
		Description: description,
		Status:      1,
		AddTime:     int(util.GetUnix()),
	}
	err := model.DB.Create(&role).Error
	if err != nil {
		c.Error(ctx, "增加部門失敗", "/role/add")
	} else {
		c.Success(ctx, "增加部門成功", "/role")
	}
}

func (c *RoleController) Edit(ctx *gin.Context) {
	RoleIntoHTML := make(map[string]interface{})
	RoleIntoHTML = IntoHTML
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/role")
		return
	}
	role := model.Role{Id: id}
	model.DB.Find(&role)
	RoleIntoHTML["role"] = role
	ctx.HTML(http.StatusOK, "role_edit.html", RoleIntoHTML)
}

func (c *RoleController) GoEdit(ctx *gin.Context) {
	title := strings.Trim(ctx.PostForm("title"), "")
	description := strings.Trim(ctx.PostForm("description"), "")
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/role")
		return
	}
	role := model.Role{Id: id}
	model.DB.Find(&role)
	role.Title = title
	role.Description = description
	err2 := model.DB.Save(&role).Error
	if err2 != nil {
		c.Error(ctx, "修改部門失敗", "/role/edit?id="+strconv.Itoa(id))
	} else {
		c.Success(ctx, "修改部門成功", "/role")
	}
}

func (c *RoleController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/role")
		return
	}
	role := model.Role{Id: id}
	administrator := []model.Administrator{}
	roleAuth := model.RoleAuth{}
	model.DB.Where("role_id=?", id).Delete(&roleAuth)
	model.DB.Preload("Role").Where("role_id=?", id).Find(&administrator)
	if len(administrator) > 0 {
		c.Error(ctx, "該部門尚有員工，無法刪除", "/role")
		return
	}
	model.DB.Delete(&role)
	c.Success(ctx, "刪除部門成功", "/role")
}

func (c *RoleController) Auth(ctx *gin.Context) {
	RoleIntoHTML := make(map[string]interface{})
	RoleIntoHTML = IntoHTML
	roleId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/role")
		return
	}
	auth := []model.Auth{}
	model.DB.Preload("AuthItem").Where("module_id=0").Find(&auth)
	roleAuth := []model.RoleAuth{}
	model.DB.Where("role_id=?", roleId).Find(&roleAuth)
	roleAuthMap := make(map[int]int)
	for _, v := range roleAuth {
		roleAuthMap[v.AuthId] = v.AuthId
	}
	for i := 0; i < len(auth); i++ {
		if _, ok := roleAuthMap[auth[i].Id]; ok {
			auth[i].Checked = true
		}
		for j := 0; j < len(auth[i].AuthItem); j++ {
			if _, ok := roleAuthMap[auth[i].AuthItem[j].Id]; ok {
				auth[i].AuthItem[j].Checked = true
			}
		}
	}
	RoleIntoHTML["authList"] = auth
	RoleIntoHTML["roleId"] = roleId
	ctx.HTML(http.StatusOK, "role_auth.html", RoleIntoHTML)
}

func (c *RoleController) GoAuth(ctx *gin.Context) {
	roleId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/role")
		return
	}
	authNode := ctx.PostFormArray("auth_node")
	roleAuth := model.RoleAuth{}
	model.DB.Where("role_id =?", roleId).Delete(&roleAuth)
	for _, v := range authNode {
		authId, _ := strconv.Atoi(v)
		roleAuth.AuthId = authId
		roleAuth.RoleId = roleId
		model.DB.Create(&roleAuth)
	}
	c.Success(ctx, "授權成功", "/role/auth?id=?"+strconv.Itoa(roleId))
}
