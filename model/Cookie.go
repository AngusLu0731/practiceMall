package model

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"practiceMall/config"
)

type cookie struct{}

var Cookie = &cookie{}

func (c cookie) Set(ctx *gin.Context, key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	ctx.SetCookie(key, string(bytes), 0, "/", config.Conf.Domain, false, false)
}

func (c cookie) Get(ctx *gin.Context, key string, obj interface{}) error {
	data, err := ctx.Cookie(key)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(data), obj)
	return err
}

func (c cookie) Delete(ctx *gin.Context, key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	ctx.SetCookie(key, string(bytes), -1, "/", config.Conf.Domain, false, false)
}
