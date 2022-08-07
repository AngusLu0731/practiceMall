package frontend

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"practiceMall/model"
	"practiceMall/util"
	"strconv"
	"strings"
)

type ProductController struct {
	BaseController
}

func NewProductController() *ProductController {
	return &ProductController{}
}

func (c *ProductController) CategoryList(ctx *gin.Context) {
	c.BaseInit(ctx)
	productIntoHTML := make(map[string]interface{})
	productIntoHTML = IntoHTML

	cateId, _ := strconv.Atoi(ctx.Param("id"))
	currentProductCate := model.ProductCate{}
	subProductCate := []model.ProductCate{}
	model.DB.Where("id = ?", cateId).Find(&currentProductCate)
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	pageSize := 5
	var tempSlice []int
	if currentProductCate.Pid == 0 {
		model.DB.Where("pid=?", currentProductCate.Id).Find(&subProductCate)
		for i := 0; i < len(subProductCate); i++ {
			tempSlice = append(tempSlice, subProductCate[i].Id)
		}
	} else {
		model.DB.Where("pid=?", currentProductCate.Pid).Find(&subProductCate)
	}
	tempSlice = append(tempSlice, cateId)
	product := []model.Product{}
	model.DB.Where("cate_id in (?)", tempSlice).Select("id,title,price,product_img,sub_title").Offset(page - 1).Order("sort desc").Find(&product)
	var count int64
	model.DB.Where("cate_id in (?)", tempSlice).Table("product").Count(&count)

	productIntoHTML["productList"] = product
	productIntoHTML["subProductCate"] = subProductCate
	productIntoHTML["currentProductCate"] = currentProductCate
	productIntoHTML["totalPages"] = math.Ceil(float64(count) / float64(pageSize))
	productIntoHTML["page"] = page

	//指定分類模板
	tpl := currentProductCate.Template
	if tpl == "" {
		tpl = "list.html"
	}
	ctx.HTML(http.StatusOK, tpl, productIntoHTML)
}

func (c *ProductController) ProductItem(ctx *gin.Context) {
	c.BaseInit(ctx)
	productIntoHTML := make(map[string]interface{})
	productIntoHTML = IntoHTML

	id := ctx.Param(":id")
	//獲取商品訊息
	product := model.Product{}
	model.DB.Where("id=?", id).Find(&product)
	productIntoHTML["product"] = product

	//獲取關聯產品  RelationProduct
	relationProduct := []model.Product{}
	product.RelationProduct = strings.ReplaceAll(product.RelationProduct, "，", ",")
	relationIds := strings.Split(product.RelationProduct, ",")
	model.DB.Where("id in (?)", relationIds).Select("id,title,price,product_version").Find(&relationProduct)
	productIntoHTML["relationProduct"] = relationProduct

	//獲取關聯贈品 ProductGift
	productGift := []model.Product{}
	product.ProductGift = strings.ReplaceAll(product.ProductGift, "，", ",")
	giftIds := strings.Split(product.ProductGift, ",")
	model.DB.Where("id in (?)", giftIds).Select("id,title,price,product_img").Find(&productGift)
	productIntoHTML["productGift"] = productGift

	//獲取關聯顏色 ProductColor
	productColor := []model.ProductColor{}
	product.ProductColor = strings.ReplaceAll(product.ProductColor, "，", ",")
	colorIds := strings.Split(product.ProductColor, ",")
	model.DB.Where("id in (?)", colorIds).Find(&productColor)
	productIntoHTML["productColor"] = productColor

	//獲取關聯配件 ProductFitting
	productFitting := []model.Product{}
	product.ProductFitting = strings.ReplaceAll(product.ProductFitting, "，", ",")
	fittingIds := strings.Split(product.ProductFitting, ",")
	model.DB.Where("id in (?)", fittingIds).Select("id,title,price,product_img").Find(&productFitting)
	productIntoHTML["productFitting"] = productFitting

	//獲取商品關聯圖片 ProductImage
	productImage := []model.ProductImage{}
	model.DB.Where("product_id=?", product.Id).Find(&productImage)
	productIntoHTML["productImage"] = productImage

	//獲取規格參數 ProductAttr
	productAttr := []model.ProductAttr{}
	model.DB.Where("product_id=?", product.Id).Find(&productAttr)
	productIntoHTML["productAttr"] = productAttr

	ctx.HTML(http.StatusOK, "item.html", productIntoHTML)
}

// Collect 收藏功能
func (c ProductController) Collect(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "參數錯誤",
		})
	}
	pId, _ := strconv.Atoi(id)
	user := model.User{}
	err := model.Cookie.Get(ctx, "userinfo", &user)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "尚未登入",
		})
	}
	isExist := model.DB.First(&user)
	if isExist.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "非法用戶",
		})
	}
	goodCollect := model.ProductCollect{}
	isExist = model.DB.Where("user_id=? AND product_id=?", user.Id, pId).First(&goodCollect)
	if isExist.RowsAffected == 0 {
		goodCollect.UserId = user.Id
		goodCollect.ProductId = pId
		goodCollect.AddTime = util.GetDate()
		model.DB.Create(&goodCollect)
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "收藏成功",
		})
	} else {
		model.DB.Delete(&goodCollect)
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "取消收藏成功",
		})
	}
}

func (c *ProductController) GetImgList(ctx *gin.Context) {
	colorId, _ := strconv.Atoi(ctx.Query("color_id"))
	productId, _ := strconv.Atoi(ctx.Query("product_id"))
	productImg := []model.ProductImage{}
	err := model.DB.Where("color_id=? AND product_id=?", colorId, productId).Find(&productImg).Error
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"result":  "失敗",
			"success": false,
		})
	} else {
		if len(productImg) == 0 {
			model.DB.Where("product_id=?", productId).Find(&productImg)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"result":  productImg,
			"success": true,
		})
	}
}
