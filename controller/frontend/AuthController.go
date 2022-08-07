package frontend

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
	"practiceMall/util"
	"regexp"
	"strings"
)

type AuthController struct {
	BaseController
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) Login(ctx *gin.Context) {
	authIntoHTML := make(map[string]interface{})
	authIntoHTML = IntoHTML
	authIntoHTML["prevPage"] = ctx.Request.Referer()
	ctx.HTML(http.StatusOK, "login.html", authIntoHTML)
}

// GoLogin 登入
func (c *AuthController) GoLogin(ctx *gin.Context) {
	phone := ctx.PostForm("phone")
	password := ctx.PostForm("password")
	phoneCode := ctx.PostForm("phone_code")
	identifyFlag := model.CaptchaVerify(ctx, phoneCode)
	userinfo := model.User{}
	model.Cookie.Get(ctx, "userinfo", &userinfo)
	if len(userinfo.Phone) == 10 {
		ctx.Redirect(http.StatusFound, "/")
	}
	if !identifyFlag {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "驗證碼錯誤或已過期",
		})
		return
	}
	password = util.Md5(password)
	user := []model.User{}
	model.DB.Where("phone=? AND password=?", phone, password).Find(&user)
	if len(user) > 0 {
		model.Cookie.Set(ctx, "userinfo", user[0])
		session := sessions.Default(ctx)
		session.Set("phone_code", phoneCode)
		session.Save()
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "登入成功",
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "帳號或密碼錯誤",
		})
		return
	}
}

// LoginOut 登出
func (c *AuthController) LoginOut(ctx *gin.Context) {
	model.Cookie.Set(ctx, "userinfo", "")
	ctx.Redirect(http.StatusFound, ctx.Request.Referer())
}

// RegisterStep1 註冊第一步
func (c *AuthController) RegisterStep1(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register_step1.html", nil)
}

// 註冊第二步
func (c *AuthController) RegisterStep2(ctx *gin.Context) {
	authIntoHTML := make(map[string]interface{})
	authIntoHTML = IntoHTML
	sign := ctx.Query("sign")
	phoneCode := ctx.Query("phone_code")
	session := sessions.Default(ctx)
	sessionPhoneCode := session.Get("phone_code")
	//確認驗證碼是否與前面的一樣
	if phoneCode != sessionPhoneCode {
		ctx.Redirect(http.StatusFound, "/auth/registerStep1")
		return
	}
	userTemp := []model.UserSms{}
	model.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		authIntoHTML["sign"] = sign
		authIntoHTML["phone_code"] = phoneCode
		authIntoHTML["phone"] = userTemp[0].Phone
		ctx.HTML(http.StatusOK, "register_step2.html", authIntoHTML)
	} else {
		ctx.Redirect(http.StatusFound, "/auth/registerStep1")
	}
}

func (c *AuthController) RegisterStep3(ctx *gin.Context) {
	authIntoHTML := make(map[string]interface{})
	authIntoHTML = IntoHTML
	sign := ctx.Query("sign")
	smsCode := ctx.Query("sms_code")
	session := sessions.Default(ctx)
	sessionSmsCode := session.Get("sms_code")
	if smsCode != sessionSmsCode {
		ctx.Redirect(http.StatusFound, "/auth/registerStep1")
		return
	}
	userTemp := []model.UserSms{}
	model.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		authIntoHTML["sign"] = sign
		authIntoHTML["sms_code"] = smsCode
		ctx.HTML(http.StatusOK, "register_step3.html", authIntoHTML)
	} else {
		ctx.Redirect(http.StatusFound, "/auth/registerStep1")
		return
	}
}

// SendCode 發送驗證碼
func (c *AuthController) SendCode(ctx *gin.Context) {
	phone := ctx.Query("phone")
	phoneCode := ctx.Query("phone_code")
	session := sessions.Default(ctx)

	identifyFlag := model.CaptchaVerify(ctx, phoneCode)
	if !identifyFlag {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "驗證碼錯誤或已過期",
		})
		return
	}
	fmt.Println(phoneCode)
	session.Set("phone_code", phoneCode)
	pattern := `^[\d]{10}$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(phone) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "手機格式錯誤",
		})
		return
	}
	user := []model.User{}
	model.DB.Where("phone = ?", phone).Find(&user)
	if len(user) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "用戶已存在",
		})
		return
	}
	addDay := util.FormatDay()
	ip := strings.Split(ctx.Request.RemoteAddr, ":")[0]
	sign := util.Md5(phone + addDay) //簽名
	smsCode := "5259"                //先固定一個值
	userTemp := []model.UserSms{}
	model.DB.Where("add_day = ? AND phone = ?", addDay, phone).Find(&userTemp)
	var sendCount int64
	model.DB.Where("add_day = ? AND ip = ?", addDay, ip).Table("user_sms").Count(&sendCount)
	if sendCount <= 10 { //判斷有無達ip當日傳送上限
		if len(userTemp) > 0 { //判斷該手機號是否在當日驗證過
			if userTemp[0].SendCount < 5 { //判斷有無達手機號當日傳送上限
				util.SendMsg(smsCode)
				session.Set("sms_code", smsCode)
				session.Save()
				oneUserSms := model.UserSms{}
				model.DB.Where("id = ?", userTemp[0].Id).Find(&oneUserSms)
				oneUserSms.SendCount += 1
				model.DB.Save(&oneUserSms)
				ctx.JSON(http.StatusOK, gin.H{
					"success":  true,
					"msg":      "驗證碼發送成功",
					"sign":     sign,
					"sms_code": smsCode,
				})
				return
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"success": false,
					"msg":     "該手機號已達當天發送次數上限",
				})
				return
			}
		} else {
			util.SendMsg(smsCode)
			session.Set("sms_code", smsCode)
			session.Save()
			oneSmsUser := model.UserSms{
				Ip:        ip,
				Phone:     phone,
				SendCount: 1,
				AddDay:    addDay,
				AddTime:   int(util.GetUnix()),
				Sign:      sign,
			}
			model.DB.Create(&oneSmsUser)
			ctx.JSON(http.StatusOK, gin.H{
				"success":  true,
				"msg":      "驗證碼發送成功",
				"sign":     sign,
				"sms_code": smsCode,
			})
			return
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "該ip已達當天發送次數上限",
		})
		return
	}
}

// ValidateSmsCode 驗證驗證碼
func (c *AuthController) ValidateSmsCode(ctx *gin.Context) {
	sign := ctx.Query("sign")
	smsCode := ctx.Query("sms_code")
	session := sessions.Default(ctx)

	userTemp := []model.UserSms{}
	model.DB.Where("sign = ?", sign).Find(&userTemp)
	if len(userTemp) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "參數錯誤",
		})
		return
	}
	sessionSmsCode := session.Get("sms_code")
	if smsCode != sessionSmsCode {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "驗證碼輸入錯誤",
		})
		return
	}
	nowTime := util.GetUnix()
	if (nowTime-int64(userTemp[0].AddTime))/1000/60 > 15 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "驗證碼過期",
		})
		return
	}
	//成功驗證
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "驗證成功",
	})
	return
}

//註冊操作
func (c *AuthController) GoRegister(ctx *gin.Context) {
	sign := ctx.PostForm("sign")
	sms_code := ctx.PostForm("sms_code")
	password := ctx.PostForm("password")
	rpassword := ctx.PostForm("rpassword")
	session := sessions.Default(ctx)
	sessionSmsCode := session.Get("sms_code")

	if sms_code != sessionSmsCode {
		ctx.Redirect(http.StatusFound, "/auth/registerStep1")
	}
	if len(password) < 6 {
		ctx.Redirect(http.StatusFound, "/auth/registerStep1")
	}
	if password != rpassword {
		ctx.Redirect(http.StatusFound, "/auth/registerStep1")
	}

	userTemp := []model.UserSms{}
	model.DB.Where("sign =?", sign).Find(&userTemp)
	ip := strings.Split(ctx.Request.RemoteAddr, ":")[0]
	if len(userTemp) > 0 {
		user := model.User{
			Phone:    userTemp[0].Phone,
			Password: util.Md5(password),
			LastIp:   ip,
		}
		model.DB.Create(&user)
		ctx.Redirect(http.StatusFound, "/")
	} else {
		ctx.Redirect(http.StatusFound, "/auth/registerStep1")
	}
}
