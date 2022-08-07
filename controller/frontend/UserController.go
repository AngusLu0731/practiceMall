package frontend

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"practiceMall/model"
	"strconv"
	"time"
)

type UserController struct {
	BaseController
}

func NewUserController() *UserController {
	return &UserController{}
}
func (c *UserController) Get(ctx *gin.Context) {
	c.BaseInit(ctx)
	userIntoHTML := make(map[string]interface{})
	userIntoHTML = IntoHTML
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	time := time.Now().Hour()
	if time >= 12 && time <= 18 {
		userIntoHTML["Hello"] = "親愛的" + user.Phone + "午安"
	} else if time >= 6 && time < 12 {
		userIntoHTML["Hello"] = "親愛的" + user.Phone + "早安"
	} else {
		userIntoHTML["Hello"] = "親愛的" + user.Phone + "晚安"
	}
	order := []model.Order{}
	model.DB.Where("uid=?", user.Id).Find(&order)
	var waitPay, waitRec int
	for i := 0; i < len(order); i++ {
		if order[i].PayStatus == 0 {
			waitPay += 1
		}
		if order[i].OrderStatus >= 2 && order[i].OrderStatus < 4 {
			waitRec += 1
		}
	}
	userIntoHTML["wait_pay"] = waitPay
	userIntoHTML["wait_rec"] = waitRec
	ctx.HTML(http.StatusOK, "welcome.html", userIntoHTML)
}
func (c *UserController) OrderList(ctx *gin.Context) {
	c.BaseInit(ctx)
	userIntoHTML := make(map[string]interface{})
	userIntoHTML = IntoHTML
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	//1.獲取訂單訊息並分頁
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	pageSize := 2 //每頁數量
	//2.獲取搜索關鍵字
	where := "uid=?"
	keywords := ctx.Query("keywords")
	if keywords != "" {
		orderItem := []model.OrderItem{}
		model.DB.Where("product_title like ?", "%"+keywords+"%").Find(&orderItem)
		var str string
		for i := 0; i < len(orderItem); i++ {
			if i == 0 {
				str += strconv.Itoa(orderItem[i].OrderId)
			} else {
				str += "," + strconv.Itoa(orderItem[i].OrderId)
			}
		}
		where += " AND id in (" + str + ")"
	}
	//3.獲取篩選條件
	orderStatus, err := strconv.Atoi(ctx.Query("order_status"))
	if err == nil {
		where += " AND order_status=" + strconv.Itoa(orderStatus)
		userIntoHTML["orderStatus"] = orderStatus
	} else {
		userIntoHTML["orderStatus"] = nil
	}
	//4.總數量
	var count int64
	model.DB.Where(where, user.Id).Table("order").Count(&count)
	order := []model.Order{}
	model.DB.Where(where, user.Id).Offset((page - 1) * pageSize).Limit(pageSize).Preload("OrderItem").Order("add_time desc").Find(&order)
	userIntoHTML["order"] = order
	userIntoHTML["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	userIntoHTML["page"] = page
	userIntoHTML["keywords"] = keywords
	ctx.HTML(http.StatusOK, "order.html", userIntoHTML)
}

func (c *UserController) OrderInfo(ctx *gin.Context) {
	c.BaseInit(ctx)
	userIntoHTML := make(map[string]interface{})
	userIntoHTML = IntoHTML
	id, _ := strconv.Atoi(ctx.Query("id"))
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	order := model.Order{}
	model.DB.Where("uid=? AND id=?", user.Id, id).Preload("order_item").Find(&order)
	userIntoHTML["order"] = order
	if order.OrderId == "" {
		ctx.Redirect(http.StatusFound, "/")
	}
	ctx.HTML(http.StatusOK, "order_info.html", userIntoHTML)
}
