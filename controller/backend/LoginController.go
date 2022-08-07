package backend

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
	"practiceMall/util"
	"strings"
)

type LoginController struct {
	BaseController
}

func NewLoginController() *LoginController {
	return &LoginController{}
}

func (c *LoginController) Get(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "backend_login.html", nil)
}

func (c *LoginController) GoLogin(ctx *gin.Context) {
	cpt := ctx.PostForm("captcha")
	flag := model.CaptchaVerify(ctx, cpt)
	if flag {
		username := strings.Trim(ctx.PostForm("username"), "")
		password := util.Md5(strings.Trim(ctx.PostForm("password"), ""))
		administrator := []model.Administrator{}
		model.DB.Where("username=? AND password=? AND status=1", username, password).Find(&administrator)
		if len(administrator) == 1 {
			session := sessions.Default(ctx)
			session.Set("userinfo", administrator[0])
			fmt.Println(administrator[0])
			err := session.Save()
			if err != nil {
				fmt.Println(err.Error())
			}
			c.Success(ctx, "登入成功", "/")
		} else {
			c.Error(ctx, "無權限或帳號密碼錯誤", "/login")
		}
	} else {
		c.Error(ctx, "驗證碼過期或錯誤", "/login")
	}
}
func (c *LoginController) LoginOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("userinfo")
	c.Success(ctx, "登出成功", "/login")
}
