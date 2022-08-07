package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
)

func FrontendAuth(ctx *gin.Context) {
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	if len(user.Phone) != 10 {
		ctx.Redirect(http.StatusFound, "/auth/login")
	}
}
