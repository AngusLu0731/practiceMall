package backend

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"practiceMall/model"
	"strconv"
)

type MainController struct {
	BaseController
}

func NewMainController() *MainController {
	return &MainController{}
}

func (c *MainController) Get(ctx *gin.Context) {
	MainIntoHTML := make(map[string]interface{})
	MainIntoHTML = IntoHTML
	session := sessions.Default(ctx)
	userinfo := session.Get("userinfo").(model.Administrator)
	if userinfo.Username != "" {
		MainIntoHTML["username"] = string(userinfo.Username)
		roleId := userinfo.RoleId
		auth := []model.Auth{}
		//獲取當前部門的權限，並把權限ID放在MAP對象中
		model.DB.Preload("AuthItem", func(db *gorm.DB) *gorm.DB {
			return db.Order("auth.sort DESC")
		}).Order("sort desc").Where("module_id=?", 0).Find(&auth)
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
		MainIntoHTML["authList"] = auth
		MainIntoHTML["isSuper"] = userinfo.IsSuper
	}
	ctx.HTML(http.StatusOK, "backend_index.html", MainIntoHTML)
}

func (c *MainController) Welcome(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "backend_welcome.html", nil)
}

func (c *MainController) ChangeStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "非法請求",
		})
	}
	table := ctx.Query("table")
	field := ctx.Query("field")
	err1 := model.DB.Exec("update "+table+" set "+field+"=ABS("+field+"-1) where id=?", id)
	if err1 != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "更新數據失敗",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "更新數據成功",
	})
}

func (c *MainController) EditNum(ctx *gin.Context) {
	id := ctx.Query("id")
	table := ctx.Query("table")
	field := ctx.Query("field")
	num := ctx.Query("num")
	err := model.DB.Exec("update " + table + " set " + field + "=" + num + " where id=" + id).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "修改數量失敗",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "修改數量成功",
	})
}
