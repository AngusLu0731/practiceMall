package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"practiceMall/model"
	"practiceMall/util"
	"strconv"
)

type MenuController struct {
	BaseController
}

func NewMenuController() *MenuController {
	return &MenuController{}
}

func (c *MenuController) Get(ctx *gin.Context) {
	MenuIntoHTML := make(map[string]interface{})
	MenuIntoHTML = IntoHTML
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	if page == 0 {
		page = 1
	}
	pageSize := 3
	menu := []model.Menu{}
	model.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&menu)
	var count int64
	model.DB.Table("menu").Count(&count)
	MenuIntoHTML["menuList"] = menu
	MenuIntoHTML["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	MenuIntoHTML["page"] = page
	ctx.HTML(http.StatusOK, "menu_index.html", MenuIntoHTML)
}

func (c *MenuController) Add(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "menu_add.html", nil)
}

func (c *MenuController) GoAdd(ctx *gin.Context) {
	title := ctx.PostForm("title")
	link := ctx.PostForm("link")
	relation := ctx.PostForm("relation")
	position, err1 := strconv.Atoi(ctx.PostForm("position"))
	isOpennew, err2 := strconv.Atoi(ctx.PostForm("isOpennew"))
	sort, err3 := strconv.Atoi(ctx.PostForm("sort"))
	status, err4 := strconv.Atoi(ctx.PostForm("status"))

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Error(ctx, "參數錯誤", "/menu/add")
		return
	}
	menu := model.Menu{
		Title:     title,
		Link:      link,
		Position:  position,
		IsOpennew: isOpennew,
		Relation:  relation,
		Sort:      sort,
		Status:    status,
		AddTime:   int(util.GetUnix()),
	}
	err := model.DB.Create(&menu).Error
	if err != nil {
		c.Error(ctx, "增加Menu失敗", "/menu/add")
		return
	} else {
		c.Success(ctx, "增加Menu成功", "/menu/add")
	}
}

func (c *MenuController) Edit(ctx *gin.Context) {
	MenuIntoHTML := make(map[string]interface{})
	MenuIntoHTML = IntoHTML
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/menu")
		return
	}
	menu := model.Menu{Id: id}
	model.DB.Find(&menu)
	MenuIntoHTML["menu"] = menu
	MenuIntoHTML["prevPage"] = ctx.Request.Referer()
	ctx.HTML(http.StatusOK, "menu_edit.html", MenuIntoHTML)
}

func (c *MenuController) GoEdit(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/menu")
		return
	}
	title := ctx.PostForm("title")
	link := ctx.PostForm("link")
	relation := ctx.PostForm("relation")
	prevPage := ctx.PostForm("prevPage")
	position, err1 := strconv.Atoi(ctx.PostForm("position"))
	isOpennew, err2 := strconv.Atoi(ctx.PostForm("isOpennew"))
	sort, err3 := strconv.Atoi(ctx.PostForm("sort"))
	status, err4 := strconv.Atoi(ctx.PostForm("status"))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.Error(ctx, "參數錯誤", "/menu/add")
		return
	}
	fmt.Println("-----------------------", relation)
	menu := model.Menu{Id: id}
	model.DB.Find(&menu)
	menu.Title = title
	menu.Link = link
	menu.Position = position
	menu.IsOpennew = isOpennew
	menu.Relation = relation
	menu.Sort = sort
	menu.Status = status

	err5 := model.DB.Save(&menu).Error
	if err5 != nil {
		c.Error(ctx, "修改數據失敗", "/menu/edit?id="+strconv.Itoa(id))
	} else {
		c.Success(ctx, "修改數據成功", prevPage)
	}
}

func (c *MenuController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/menu")
		return
	}
	menu := model.Menu{Id: id}
	model.DB.Delete(&menu)
	c.Success(ctx, "刪除數據成功", ctx.Request.Referer())
}
