package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
	"strconv"
)

type CartController struct {
	BaseController
}

func NewCartController() *CartController {
	return &CartController{}
}

func (c *CartController) Get(ctx *gin.Context) {
	c.BaseInit(ctx)
	cartIntoHTML := make(map[string]interface{})
	cartIntoHTML = IntoHTML

	cartList := []model.Cart{}
	model.Cookie.Get(ctx, "cartList", &cartList)

	var allPrice float64
	//計算總價
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	cartIntoHTML["cartList"] = cartList
	cartIntoHTML["allPrice"] = allPrice
	ctx.HTML(http.StatusOK, "cart.html", cartIntoHTML)
}
func (c *CartController) AddCart(ctx *gin.Context) {
	c.BaseInit(ctx)
	cartIntoHTML := make(map[string]interface{})
	cartIntoHTML = IntoHTML

	colorId, _ := strconv.Atoi(ctx.Query("color_id"))
	productId, _ := strconv.Atoi(ctx.Query("product_id"))
	product := model.Product{}
	productColor := model.ProductColor{}
	err1 := model.DB.Where("id=?", productId).Find(&product).Error
	err2 := model.DB.Where("id=?", colorId).Find(&productColor).Error

	if err1 != nil || err2 != nil {
		ctx.Redirect(http.StatusFound, "/item_"+strconv.Itoa(product.Id)+".html")
		return
	}

	//1.獲取增加購物車的商品資訊
	currentData := model.Cart{
		Id:             productId,
		Title:          product.Title,
		Price:          product.Price,
		ProductVersion: product.ProductVersion,
		Num:            1,
		ProductColor:   productColor.ColorName,
		ProductImg:     product.ProductImg,
		ProductGift:    product.ProductGift,
		ProductAttr:    "",
		Checked:        true, //默認勾選
	}

	//2.判斷購物車內有無數據(Cookie)
	cartList := []model.Cart{}
	model.Cookie.Get(ctx, "cartList", &cartList)
	if len(cartList) > 0 { //購物車內有數據
		//3.判斷購物車內有無當前數據
		if model.CartHasData(cartList, currentData) {
			for i := 0; i < len(cartList); i++ {
				if cartList[i].Id == currentData.Id && cartList[i].ProductColor == currentData.ProductColor && cartList[i].ProductAttr == currentData.ProductAttr {
					cartList[i].Num = cartList[i].Num + 1
				}
			}
		} else {
			cartList = append(cartList, currentData)
		}
		model.Cookie.Set(ctx, "cartList", cartList)
	} else {
		//4.如購物車內沒有任何數據，直接將當前數據寫入cookie
		cartList = append(cartList, currentData)
		model.Cookie.Set(ctx, "cartList", cartList)
	}
	cartIntoHTML["product"] = product
	ctx.HTML(http.StatusOK, "addcart_success.html", cartIntoHTML)
}

// DecCart 減少商品數量
func (c *CartController) DecCart(ctx *gin.Context) {
	var (
		flag            bool
		allPrice        float64
		currentAllPrice float64
		num             int
	)
	productId, _ := strconv.Atoi(ctx.Query("color_id"))
	productColor := ctx.Query("product_color")
	productAttr := ""

	cartList := []model.Cart{}
	model.Cookie.Get(ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId && cartList[i].ProductColor == productColor && cartList[i].ProductAttr == productAttr {
			if cartList[i].Num > 1 {
				cartList[i].Num = cartList[i].Num - 1
			}
			flag = true
			num = cartList[i].Num
			currentAllPrice = cartList[i].Price * float64(cartList[i].Num)
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	if flag {
		model.Cookie.Set(ctx, "cartList", cartList)
		ctx.JSON(http.StatusOK, gin.H{
			"success":         true,
			"message":         "修改數量成功",
			"allPrice":        allPrice,
			"currentAllPrice": currentAllPrice,
			"num":             num,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "參數錯誤",
		})
	}
}

// IncCart 增加商品數量
func (c *CartController) IncCart(ctx *gin.Context) {
	var (
		flag            bool
		allPrice        float64
		currentAllPrice float64
		num             int
	)
	productId, _ := strconv.Atoi(ctx.Query("color_id"))
	productColor := ctx.Query("product_color")
	productAttr := ""

	cartList := []model.Cart{}
	model.Cookie.Get(ctx, "cartList", &cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId && cartList[i].ProductColor == productColor && cartList[i].ProductAttr == productAttr {
			cartList[i].Num = cartList[i].Num + 1
			flag = true
			num = cartList[i].Num
			currentAllPrice = cartList[i].Price * float64(cartList[i].Num)
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	if flag {
		model.Cookie.Set(ctx, "cartList", cartList)
		ctx.JSON(http.StatusOK, gin.H{
			"success":         true,
			"message":         "修改數量成功",
			"allPrice":        allPrice,
			"currentAllPrice": currentAllPrice,
			"num":             num,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "參數錯誤",
		})
	}
}

// ChangeOneCart 勾選或取消勾選
func (c *CartController) ChangeOneCart(ctx *gin.Context) {
	var (
		flag     bool
		allPrice float64
	)
	productId, _ := strconv.Atoi(ctx.Query("color_id"))
	productColor := ctx.Query("product_color")
	productAttr := ""

	cartList := []model.Cart{}
	model.Cookie.Get(ctx, "cartList", cartList)

	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId && cartList[i].ProductColor == productColor && cartList[i].ProductAttr == productAttr {
			cartList[i].Checked = !cartList[i].Checked
			flag = true
		}
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	if flag {
		model.Cookie.Set(ctx, "cartList", cartList)
		ctx.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "修改狀態成功",
			"allPrice": allPrice,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "傳入參數錯誤",
		})
	}
}

// ChangeAllCart 全選反選
func (c *CartController) ChangeAllCart(ctx *gin.Context) {
	flag, _ := strconv.Atoi(ctx.Query("flag"))
	var allPrice float64
	cartList := []model.Cart{}
	model.Cookie.Get(ctx, "cartList", cartList)

	for i := 0; i < len(cartList); i++ {
		if flag == 1 {
			cartList[i].Checked = true
		} else {
			cartList[i].Checked = false
		}
		//計算總價
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
		}
	}
	model.Cookie.Set(ctx, "cartList", cartList)
	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"allPrice": allPrice,
	})
}

func (c *CartController) DelCart(ctx *gin.Context) {
	productId, _ := strconv.Atoi(ctx.Query("color_id"))
	productColor := ctx.Query("product_color")
	productAttr := ""
	cartList := []model.Cart{}
	model.Cookie.Get(ctx, "cartList", cartList)
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Id == productId && cartList[i].ProductColor == productColor && cartList[i].ProductAttr == productAttr {
			//執行刪除
			cartList = append(cartList[:i], cartList[(i+1):]...)
		}
	}
	model.Cookie.Set(ctx, "cartList", cartList)
	ctx.Redirect(http.StatusFound, "/cart")
}
