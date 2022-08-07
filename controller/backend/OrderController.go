package backend

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"practiceMall/model"
	"strconv"
)

type OrderController struct {
	BaseController
}

func NewOrderController() *OrderController {
	return &OrderController{}
}

func (c *OrderController) Get(ctx *gin.Context) {
	OrderIntoHTML := make(map[string]interface{})
	OrderIntoHTML = IntoHTML
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page == 0 {
		page = 1
	}
	pageSize := 5
	keyword := ctx.Query("keyword")
	order := []model.Order{}
	var count int64
	if keyword == "" {
		model.DB.Table("order").Count(&count)
		model.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&order)
	} else {
		model.DB.Where("phone=?", keyword).Offset((page - 1) * pageSize).Limit(pageSize).Find(&order)
		model.DB.Where("phone=?", keyword).Table("order").Count(&count)
	}
	OrderIntoHTML["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	OrderIntoHTML["page"] = page
	OrderIntoHTML["order"] = order
	ctx.HTML(http.StatusOK, "order_order.html", OrderIntoHTML)
}

func (c *OrderController) Detail(ctx *gin.Context) {
	ctx.String(http.StatusOK, "詳情頁面")
}

func (c *OrderController) Edit(ctx *gin.Context) {
	OrderIntoHTML := make(map[string]interface{})
	OrderIntoHTML = IntoHTML
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "order")
		return
	}
	order := model.Order{}
	model.DB.Where("id = ?", id).Find(&order)
	OrderIntoHTML["order"] = order
	ctx.HTML(http.StatusOK, "order_edit.html", OrderIntoHTML)
}

func (c *OrderController) GoEdit(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "order")
		return
	}
	orderId := ctx.PostForm("order_id")
	allPrice := ctx.PostForm("all_price")
	name := ctx.PostForm("name")
	phone := ctx.PostForm("phone")
	address := ctx.PostForm("address")
	zipcode := ctx.PostForm("zipcode")
	payStatus, _ := strconv.Atoi(ctx.PostForm("pay_status"))
	payType, _ := strconv.Atoi(ctx.PostForm("pay_type"))
	orderStatus, _ := strconv.Atoi(ctx.PostForm("order_status"))
	order := model.Order{}
	model.DB.Where("id=?", id).Find(&order)
	order.OrderId = orderId
	order.AllPrice, _ = strconv.ParseFloat(allPrice, 64)
	order.Name = name
	order.Phone = phone
	order.Address = address
	order.Zipcode = zipcode
	order.PayStatus = payStatus
	order.PayType = payType
	order.OrderStatus = orderStatus
	model.DB.Save(&order)
	c.Success(ctx, "訂單修改成功", "/order")
}
func (c *OrderController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "order")
		return
	}
	order := model.Order{}
	model.DB.Where("id=?", id).Delete(&order)
	c.Success(ctx, "刪除訂單成功", "/order")
}
