package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
	"practiceMall/util"
	"strconv"
	"strings"
)

type AdministratorController struct {
	BaseController
}

func NewAdministratorController() *AdministratorController {
	return &AdministratorController{}
}

func (c *AdministratorController) Get(ctx *gin.Context) {
	AdministratorIntoHTML := make(map[string]interface{})
	AdministratorIntoHTML = IntoHTML
	administrator := []model.Administrator{}
	model.DB.Preload("role").Find(&administrator)
	AdministratorIntoHTML["administratorList"] = administrator
	ctx.HTML(http.StatusOK, "administrator_index.html", AdministratorIntoHTML)
}

func (c *AdministratorController) Add(ctx *gin.Context) {
	AdministratorIntoHTML := make(map[string]interface{})
	AdministratorIntoHTML = IntoHTML
	role := []model.Role{}
	model.DB.Find(&role)
	AdministratorIntoHTML["roleList"] = role
	ctx.HTML(http.StatusOK, "administrator_add.html", AdministratorIntoHTML)
}

func (c *AdministratorController) GoAdd(ctx *gin.Context) {
	username := strings.Trim(ctx.PostForm("username"), "")
	password := strings.Trim(ctx.PostForm("password"), "")
	mobile := strings.Trim(ctx.PostForm("mobile"), "")
	email := strings.Trim(ctx.PostForm("email"), "")
	roleId, err1 := strconv.Atoi(ctx.PostForm("role_id"))
	if err1 != nil {
		c.Error(ctx, "非法請求", "administrator/add")
		return
	}
	if len(username) < 2 || len(password) < 6 {
		c.Error(ctx, "使用者名稱或密碼長度有誤", "administrator/add")
		return
	} else if !util.VerifyEmail(email) {
		c.Error(ctx, "電子郵件錯誤", "administrator/add")
		return
	}
	administratorList := []model.Administrator{}
	model.DB.Where("username = ?", username).Find(&administratorList)
	if len(administratorList) > 0 {
		c.Error(ctx, "使用者已存在", "administrator/add")
		return
	}
	administrator := model.Administrator{
		Username: username,
		Password: util.Md5(password),
		Mobile:   mobile,
		Email:    email,
		Status:   1,
		AddTime:  int(util.GetUnix()),
		RoleId:   roleId,
	}
	err := model.DB.Create(&administrator).Error
	if err != nil {
		c.Error(ctx, "增加管理員失敗", "/administrator/add")
		return
	}
	c.Success(ctx, "增加管理員成功", "/administrator")
}

func (c *AdministratorController) Edit(ctx *gin.Context) {
	AdministratorIntoHTML := make(map[string]interface{})
	AdministratorIntoHTML = IntoHTML
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "傳入參數錯誤", "/administrator")
		return
	}
	administrator := model.Administrator{}
	model.DB.Where("id = ?", id).Find(&administrator)
	AdministratorIntoHTML["administrator"] = administrator
	role := []model.Role{}
	model.DB.Find(&role)
	AdministratorIntoHTML["roleList"] = role
	ctx.HTML(http.StatusOK, "administrator_edit.html", AdministratorIntoHTML)
}

func (c *AdministratorController) GoEdit(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "傳入參數錯誤", "/administrator")
		return
	}
	username := strings.Trim(ctx.PostForm("username"), "")
	password := strings.Trim(ctx.PostForm("password"), "")
	mobile := strings.Trim(ctx.PostForm("mobile"), "")
	email := strings.Trim(ctx.PostForm("email"), "")
	roleId, err1 := strconv.Atoi(ctx.PostForm("role_id"))
	if err1 != nil {
		c.Error(ctx, "非法請求", "/administrator")
		return
	}
	if password != "" {
		if len(password) < 6 {
			c.Error(ctx, "密碼長度需大於6碼", "/administrator/add?id="+strconv.Itoa(id))
			return
		} else if !util.VerifyEmail(email) {
			c.Error(ctx, "電子郵件錯誤", "administrator/add?id="+strconv.Itoa(id))
			return
		}
		password = util.Md5(password)
	}
	administrator := model.Administrator{Id: id}
	model.DB.Find(&administrator)
	administrator.Username = username
	administrator.Password = password
	administrator.Mobile = mobile
	administrator.Email = email
	administrator.RoleId = roleId
	err2 := model.DB.Save(&administrator).Error
	if err2 != nil {
		c.Error(ctx, "修改管理員失敗", "/administrator/edit?id="+strconv.Itoa(id))
	} else {
		c.Success(ctx, "修改管理員成功", "/administrator")
	}
}

func (c *AdministratorController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "傳入參數錯誤", "/administrator")
		return
	}
	administrator := model.Administrator{Id: id}
	model.DB.Delete(&administrator)
	c.Success(ctx, "刪除管理員成功", "/administrator")
}
