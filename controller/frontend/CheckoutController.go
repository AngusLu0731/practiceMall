package frontend

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
	"practiceMall/util"
	"strconv"
)

type CheckoutController struct {
	BaseController
}

func NewCheckOutController() *CheckoutController {
	return &CheckoutController{}
}

func (c *CheckoutController) CheckOut(ctx *gin.Context) {
	c.BaseInit(ctx)
	checkOutIntoHTML := make(map[string]interface{})
	checkOutIntoHTML = IntoHTML
	//1.獲取要結帳的商品
	cartList := []model.Cart{}
	orderList := []model.Cart{} //要結帳商品
	model.Cookie.Get(ctx, "cartList", &cartList)

	var allPrice float64
	//2.計算總價
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
		}
	}
	//3.判斷結帳slice內有無商品
	if len(orderList) == 0 {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	checkOutIntoHTML["orderList"] = orderList
	checkOutIntoHTML["allPrice"] = allPrice

	//4.獲取收貨地址
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	addressList := []model.Address{}
	model.DB.Where("uid=?", user.Id).Order("default_address desc").Find(&addressList)
	checkOutIntoHTML["addressList"] = addressList

	//5.為了防止重複提交，生成簽名
	orderSign := util.Md5(util.GetRandomNum())
	session := sessions.Default(ctx)
	session.Set("orderSign", orderSign)
	session.Save()
	checkOutIntoHTML["orderSign"] = orderSign

	ctx.HTML(http.StatusOK, "checkout.html", checkOutIntoHTML)
}

func (c *CheckoutController) GoOrder(ctx *gin.Context) {
	//1.防止重複提交訂單
	orderSign := ctx.Query("orderSign")
	session := sessions.Default(ctx)
	sessionOrderSign := session.Get("orderSign")
	if orderSign != sessionOrderSign {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	session.Delete("orderSign")
	//2.獲取收貨地址
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	addressResult := []model.Address{}
	model.DB.Where("uid=? AND default_address=1", user.Id).Find(&addressResult)
	if len(addressResult) > 0 {
		//2.獲取要購買的商品信息
		cartList := []model.Cart{}
		orderList := []model.Cart{}
		var allPrice float64
		model.Cookie.Get(ctx, "cartList", &cartList)
		for i := 0; i < len(cartList); i++ {
			if cartList[i].Checked {
				allPrice += cartList[i].Price * float64(cartList[i].Num)
				orderList = append(orderList, cartList[i])
			}
		}
		//3.把訂單訊息放在訂單表，商品訊息放商品表
		order := model.Order{
			OrderId:     util.GenerateOrderId(),
			Uid:         user.Id,
			AllPrice:    allPrice,
			Phone:       addressResult[0].Phone,
			Name:        addressResult[0].Name,
			Address:     addressResult[0].Address,
			Zipcode:     addressResult[0].Zipcode,
			PayStatus:   0,
			PayType:     0,
			OrderStatus: 0,
			AddTime:     int(util.GetUnix()),
		}
		err := model.DB.Create(&order).Error
		if err == nil {
			for i := 0; i < len(orderList); i++ {
				orderItem := model.OrderItem{
					OrderId:        order.Id,
					Uid:            user.Id,
					ProductTitle:   orderList[i].Title,
					ProductId:      orderList[i].Id,
					ProductImg:     orderList[i].ProductImg,
					ProductPrice:   orderList[i].Price,
					ProductNum:     orderList[i].Num,
					ProductVersion: orderList[i].ProductVersion,
					ProductColor:   orderList[i].ProductColor,
					AddTime:        int(util.GetUnix()),
				}
				err := model.DB.Create(&orderItem).Error
				if err != nil {
					fmt.Println(err)
				}
			}
			//4.刪除購物車內的選中數據
			noSelectedCartList := []model.Cart{}
			for i := 0; i < len(cartList); i++ {
				if !cartList[i].Checked {
					noSelectedCartList = append(noSelectedCartList, cartList[i])
				}
			}
			model.Cookie.Set(ctx, "cartList", noSelectedCartList)
			ctx.Redirect(http.StatusFound, "/buy/confirm?id="+strconv.Itoa(order.Id))
		} else {
			//錯誤請求
			ctx.Redirect(http.StatusFound, "/")
		}
	} else {
		//錯誤請求
		ctx.Redirect(http.StatusFound, "/")
	}
}

// Confirm 確認結算
func (c *CheckoutController) Confirm(ctx *gin.Context) {
	c.BaseInit(ctx)
	checkOutIntoHTML := make(map[string]interface{})
	checkOutIntoHTML = IntoHTML
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	//獲取用戶訊息
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	//獲取訂單訊息
	order := model.Order{}
	model.DB.Where("id=?", id).Find(&order)
	checkOutIntoHTML["order"] = order
	if order.Uid != user.Id {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	//獲取訂單下的商品訊息
	orderItem := []model.OrderItem{}
	model.DB.Where("order_id=?", id).Find(&orderItem)
	checkOutIntoHTML["orderItem"] = orderItem

	ctx.HTML(http.StatusOK, "confirm.html", checkOutIntoHTML)
}

// OrderPayStatus 獲取訂單支付狀態
func (c *CheckoutController) OrderPayStatus(ctx *gin.Context) {
	//1.獲取訂單號碼
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "參數錯誤",
		})
		return
	}
	//2.查詢訂單
	order := model.Order{}
	model.DB.Where("id=?", id).Find(&order)
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	//3.判斷當前訊息是否合法
	if order.Uid != user.Id {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "參數錯誤",
		})
		return
	}
	//4.判斷訂單支付狀態
	if order.PayStatus == 1 && order.OrderStatus == 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "已支付",
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "未支付",
		})
		return
	}
}
