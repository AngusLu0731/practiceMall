package util

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"practiceMall/config"
	"practiceMall/model"
	"strings"
)

func BackendAuth(ctx *gin.Context) {
	pathname := ctx.Request.URL.String()
	session := sessions.Default(ctx)
	userinfo, ok := session.Get("userinfo").(model.Administrator)
	fmt.Println(userinfo)
	if !(ok && userinfo.Username != "") {
		if pathname != "/"+config.Conf.AdminPath+"/login" &&
			pathname != "/"+config.Conf.AdminPath+"/login/gologin" &&
			pathname != "/"+config.Conf.AdminPath+"/login/verificode" {
			ctx.Redirect(http.StatusFound, "/"+config.Conf.AdminPath+"/login")
		}
	} else {
		pathname = strings.Replace(pathname, "/"+config.Conf.AdminPath, "", 1)
		urlPath, _ := url.Parse(pathname)
		if userinfo.IsSuper == 0 && !excludeAuthPath(string(urlPath.Path)) {
			roleId := userinfo.RoleId
			roleAuth := []model.RoleAuth{}
			model.DB.Where("role_id=?", roleId).Find(&roleAuth)
			roleAuthMap := make(map[int]int)
			for _, v := range roleAuth {
				roleAuthMap[v.AuthId] = v.AuthId
			}
			auth := model.Auth{}
			model.DB.Where("url=?", urlPath.Path).Find(&auth)
			if _, ok := roleAuthMap[auth.Id]; !ok {
				ctx.String(http.StatusOK, "沒有權限")
				return
			}
		}
	}
}

func excludeAuthPath(urlPath string) bool {
	excludeAuthPathSlice := strings.Split(config.Conf.ExcludeAuthPath, ",")
	for _, v := range excludeAuthPathSlice {
		if v == urlPath {
			return true
		}
	}
	return false
}
