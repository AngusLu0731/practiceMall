package main

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"practiceMall/controller/backend"
	"practiceMall/controller/frontend"
	"practiceMall/model"
	"practiceMall/util"
	"strings"
)

func main() {
	model.ConnectMySql()
	defer model.CloseMysql()

	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	gob.Register(model.Administrator{})

	//Set HTML FuncMap
	router.SetFuncMap(template.FuncMap{
		"md5":             util.Md5,
		"getUnixNano":     util.GetUnixNano,
		"mul":             util.Mul,
		"formatImage":     util.FormatImg,
		"substr":          util.SubStr,
		"formatAttribute": util.FormatAttr,
		"timestampToDate": util.TimestampToDate,
		"setting":         model.GetSettingByColumn,
	})
	//Load HTML File
	var httpFiles []string
	filepath.Walk("./view", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			httpFiles = append(httpFiles, path)
		}
		return nil
	})
	router.LoadHTMLFiles(httpFiles...)

	//Load Static File
	router.Static("/static", "./static")

	/****************************************************
		Router Setting
	*****************************************************/

	/******************************************
		FrontEnd
	*******************************************/
	router.GET("", frontend.NewIndexController().Get)
	router.GET("/category_:id([0-9]+).html", frontend.NewProductController().CategoryList)
	router.GET("/item_:id([0-9]+).html", frontend.NewProductController().ProductItem)
	router.GET("/product/getImgList", frontend.NewProductController().GetImgList)
	router.GET("/product/collect", frontend.NewProductController().Collect)
	router.GET("/captcha", func(ctx *gin.Context) {
		model.Captcha(ctx, 4)
	})

	//購物車 Cart
	cart := router.Group("/cart")
	cart.GET("", frontend.NewCartController().Get)
	cart.GET("/addCart", frontend.NewCartController().AddCart)
	cart.GET("/incCart", frontend.NewCartController().IncCart)
	cart.GET("/decCart", frontend.NewCartController().DecCart)
	cart.GET("/delCart", frontend.NewCartController().DelCart)
	cart.GET("/changeOneCart", frontend.NewCartController().ChangeOneCart)
	cart.GET("/changeAllCart", frontend.NewCartController().ChangeAllCart)

	//驗證 Auth
	auth := router.Group("/auth")
	auth.GET("/sendCode", frontend.NewAuthController().SendCode)
	auth.POST("/doRegister", frontend.NewAuthController().GoRegister)
	auth.GET("/validateSmsCode", frontend.NewAuthController().ValidateSmsCode)
	auth.GET("/registerStep1", frontend.NewAuthController().RegisterStep1)
	auth.GET("/registerStep2", frontend.NewAuthController().RegisterStep2)
	auth.GET("/registerStep3", frontend.NewAuthController().RegisterStep3)
	auth.GET("/login", frontend.NewAuthController().Login)
	auth.POST("/goLogin", frontend.NewAuthController().GoLogin)
	auth.GET("/loginOut", frontend.NewAuthController().LoginOut)

	//結帳 Buy
	buy := router.Group("/buy").Use(util.FrontendAuth)
	buy.GET("/checkout", frontend.NewCheckOutController().CheckOut)
	buy.POST("/doOrder", frontend.NewCheckOutController().GoOrder)
	buy.GET("/confirm", frontend.NewCheckOutController().Confirm)
	buy.GET("/orderPayStatus", frontend.NewCheckOutController().OrderPayStatus)

	//會員地址 address
	address := router.Group("/address").Use(util.FrontendAuth)
	address.POST("/addAddress", frontend.NewAddressController().AddAddress)
	address.GET("/getOneAddressList", frontend.NewAddressController().GetOneAddressList)
	address.POST("/goEditAddressList", frontend.NewAddressController().GoEditAddressList)
	address.GET("/changeDefaultAddress", frontend.NewAddressController().ChangeDefaultAddress)

	//用戶 user
	user := router.Group("/user").Use(util.FrontendAuth)
	user.GET("", frontend.NewUserController().Get)
	user.GET("/order", frontend.NewUserController().OrderList)
	user.GET("/orderinfo", frontend.NewUserController().OrderInfo)

	/******************************************
		BackEnd
	*******************************************/
	admin := router.Group("/backend").Use(util.BackendAuth)
	//後臺管理
	admin.GET("/", backend.NewMainController().Get)
	admin.GET("/welcome", backend.NewMainController().Welcome)
	admin.GET("/main/changestatus", backend.NewMainController().ChangeStatus)
	admin.GET("/main/editnum", backend.NewMainController().EditNum)
	admin.GET("/login", backend.NewLoginController().Get)
	admin.POST("/login/gologin", backend.NewLoginController().GoLogin)
	admin.GET("/login/loginout", backend.NewLoginController().LoginOut)
	admin.GET("/banner", backend.NewBannerController().Get)

	//管理員管理
	administrator := router.Group("/backend/administrator")
	administrator.GET("/", backend.NewAdministratorController().Get)
	administrator.GET("/add", backend.NewAdministratorController().Add)
	administrator.GET("/edit", backend.NewAdministratorController().Edit)
	administrator.GET("/delete", backend.NewAdministratorController().Delete)
	administrator.POST("/goadd", backend.NewAdministratorController().GoAdd)
	administrator.POST("/goedit", backend.NewAdministratorController().GoEdit)

	//部門管理
	role := router.Group("/backend/role")
	role.GET("/", backend.NewRoleController().Get)
	role.GET("/add", backend.NewRoleController().Add)
	role.GET("/edit", backend.NewRoleController().Edit)
	role.GET("/auth", backend.NewRoleController().Auth)
	role.GET("/delete", backend.NewRoleController().Delete)
	role.POST("/goadd", backend.NewRoleController().GoAdd)
	role.POST("/goedit", backend.NewRoleController().GoEdit)
	role.POST("/goauth", backend.NewRoleController().GoAuth)

	//權限管理
	authBackend := router.Group("/backend/auth")
	authBackend.GET("/", backend.NewAuthController().Get)
	authBackend.GET("/add", backend.NewAuthController().Add)
	authBackend.GET("/edit", backend.NewAuthController().Edit)
	authBackend.GET("/delete", backend.NewAuthController().Delete)
	authBackend.POST("/goadd", backend.NewAuthController().GoAdd)
	authBackend.POST("/goedit", backend.NewAuthController().GoEdit)

	//商品管理(待處理)
	admin.GET("/product/*id", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "待處理")
	})
	admin.GET("/productType/*id", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "待處理")
	})
	admin.GET("/productCate/*id", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "待處理")
	})
	admin.GET("/productTypeAttribute/*id", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "待處理")
	})

	//輪播圖管理
	banner := router.Group("/backend/banner")
	banner.GET("/", backend.NewBannerController().Get)
	banner.GET("/add", backend.NewBannerController().Add)
	banner.GET("/edit", backend.NewBannerController().Edit)
	banner.GET("/delete", backend.NewBannerController().Delete)
	banner.POST("/goadd", backend.NewBannerController().GoAdd)
	banner.POST("/goedit", backend.NewBannerController().GoEdit)

	//訂單管理
	order := router.Group("/backend/order")
	order.GET("/", backend.NewOrderController().Get)
	order.GET("/detail", backend.NewOrderController().Detail)
	order.GET("/edit", backend.NewOrderController().Edit)
	order.GET("/delete", backend.NewOrderController().Delete)
	order.POST("/goedit", backend.NewOrderController().GoEdit)

	//索引管理
	menu := router.Group("/backend/menu")
	menu.GET("/", backend.NewMenuController().Get)
	menu.GET("/add", backend.NewMenuController().Add)
	menu.GET("/edit", backend.NewMenuController().Edit)
	menu.GET("/delete", backend.NewMenuController().Delete)
	menu.POST("/goadd", backend.NewMenuController().GoAdd)
	menu.POST("/goedit", backend.NewMenuController().GoEdit)

	//系統設置
	setting := router.Group("/backend/setting")
	setting.GET("/", backend.NewSettingController().Get)
	setting.POST("/goedit", backend.NewSettingController().GoEdit)
	/****************************************************
		Router Setting End
	*****************************************************/

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
