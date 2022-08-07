package model

import (
	"bytes"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var Store captcha.Store

func init() {
	Store = captcha.NewMemoryStore(60, 10*time.Minute)
	captcha.SetCustomStore(Store)
}

func Captcha(ctx *gin.Context, length ...int) {
	l := 4
	w, h := 100, 40
	if len(length) == 1 {
		l = length[0]
	}
	if len(length) == 2 {
		w = length[1]
	}
	if len(length) == 3 {
		h = length[2]
	}
	captchaId := captcha.NewLen(l)
	cId := []byte(captchaId)
	Store.Set(ctx.Request.RemoteAddr, cId)
	session := sessions.Default(ctx)
	session.Set("phone_code", cId)
	_ = Serve(ctx.Writer, ctx.Request, captchaId, ".png", "zh", false, w, h)
}
func CaptchaVerify(ctx *gin.Context, code string) bool {
	if captchaId := Store.Get(ctx.Request.RemoteAddr, true); captchaId != nil {
		if captcha.VerifyString(string(captchaId), code) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
