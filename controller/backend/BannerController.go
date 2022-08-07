package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"practiceMall/model"
	"practiceMall/util"
	"strconv"
)

type BannerController struct {
	BaseController
}

func NewBannerController() *BannerController {
	return &BannerController{}
}

func (c *BannerController) Get(ctx *gin.Context) {
	BannerIntoHTML := make(map[string]interface{})
	BannerIntoHTML = IntoHTML
	banner := []model.Banner{}
	model.DB.Find(&banner)
	BannerIntoHTML["bannerList"] = banner
	ctx.HTML(http.StatusOK, "banner_index.html", BannerIntoHTML)
}

func (c *BannerController) Add(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "banner_add.html", nil)
}

func (c *BannerController) GoAdd(ctx *gin.Context) {
	title := ctx.PostForm("title")
	link := ctx.PostForm("link")
	bannerType, err1 := strconv.Atoi(ctx.PostForm("banner_type"))
	sort, err2 := strconv.Atoi(ctx.PostForm("sort"))
	status, err3 := strconv.Atoi(ctx.PostForm("status"))
	if err1 != nil || err3 != nil {
		c.Error(ctx, "非法請求", "/banner")
		return
	}
	if err2 != nil {
		c.Error(ctx, "排序內的內容不合法", "/banner/add")
		return
	}
	bannerImgSrc, err4 := c.LocalUploadImage(ctx, "banner_img")
	if err4 == nil {
		banner := model.Banner{
			Title:      title,
			BannerType: bannerType,
			BannerImg:  bannerImgSrc,
			Link:       link,
			Sort:       sort,
			Status:     status,
			AddTime:    int(util.GetUnix()),
		}
		model.DB.Create(&banner)
		c.Success(ctx, "增加輪播圖成功", "/banner")
	} else {
		c.Error(ctx, "增加輪播圖失敗", "/banner/add")
		return
	}
}

func (c *BannerController) Edit(ctx *gin.Context) {
	BannerIntoHTML := make(map[string]interface{})
	BannerIntoHTML = IntoHTML
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		c.Error(ctx, "參數錯誤", "/banner")
		return
	}
	banner := model.Banner{Id: id}
	model.DB.Find(&banner)
	BannerIntoHTML["banner"] = banner
	ctx.HTML(http.StatusOK, "banner_edit.html", BannerIntoHTML)
}

func (c *BannerController) GoEdit(ctx *gin.Context) {
	title := ctx.PostForm("title")
	link := ctx.PostForm("link")
	id, err := strconv.Atoi(ctx.PostForm("id"))
	bannerType, err1 := strconv.Atoi(ctx.PostForm("banner_type"))
	sort, err2 := strconv.Atoi(ctx.PostForm("sort"))
	status, err3 := strconv.Atoi(ctx.PostForm("status"))
	if err != nil || err1 != nil || err3 != nil {
		c.Error(ctx, "非法請求", "/banner")
		return
	}
	if err2 != nil {
		c.Error(ctx, "排序內的內容不合法", "/banner/add")
		return
	}
	bannerImgSrc, _ := c.LocalUploadImage(ctx, "banner_img")
	banner := model.Banner{Id: id}
	model.DB.Find(&banner)
	banner.Title = title
	banner.BannerType = bannerType
	banner.Link = link
	banner.Sort = sort
	banner.Status = status
	if bannerImgSrc != "" {
		banner.BannerImg = bannerImgSrc
	}
	err5 := model.DB.Save(&banner).Error
	if err5 != nil {
		c.Error(ctx, "修改輪播圖失敗", "/banner/edit?id="+strconv.Itoa(id))
		return
	}
	c.Success(ctx, "修改輪播圖成功", "/banner")
}

func (c *BannerController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		c.Error(ctx, "傳入參數錯誤", "/banner")
		return
	}
	banner := model.Banner{Id: id}
	model.DB.Find(&banner)
	address := "C:/Users/Angus/GolandProjects/practiceMall/" + banner.BannerImg
	test := os.Remove(address)
	if test != nil {
		c.Error(ctx, "刪除圖片錯誤，請確認", "/banner")
		return
	}
	model.DB.Delete(&banner)
	c.Success(ctx, "刪除輪播圖成功", "/banner")
}
