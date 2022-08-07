package frontend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practiceMall/model"
	"strconv"
)

type AddressController struct {
	BaseController
}

func NewAddressController() *AddressController {
	return &AddressController{}
}

// AddAddress 增加地址
func (c *AddressController) AddAddress(ctx *gin.Context) {
	user := model.User{}
	name := ctx.Query("name")
	phone := ctx.Query("phone")
	address := ctx.Query("address")
	zipcode := ctx.Query("zipcode")
	var addressCount int64
	model.DB.Where("uid=?", user.Id).Table("address").Count(&addressCount)
	if addressCount > 10 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "增加地址失敗，地址數量超過上限",
		})
		return
	}
	//將其他地址的預設取消
	model.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	addressResult := model.Address{
		Uid:            user.Id,
		Phone:          phone,
		Name:           name,
		Address:        address,
		Zipcode:        zipcode,
		DefaultAddress: 1, //新增的地址設為預設
	}
	model.DB.Create(&addressResult)
	allAddressResult := []model.Address{}
	model.DB.Where("uid=?", user.Id).Find(&allAddressResult)
	ctx.JSON(http.StatusOK, gin.H{
		"success": false,
		"result":  allAddressResult,
	})
}

func (c *AddressController) GetOneAddressList(ctx *gin.Context) {
	addressId, err := strconv.Atoi(ctx.Query("address_id"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "參數錯誤",
		})
		return
	}
	address := model.Address{}
	model.DB.Where("id=?", addressId).Find(&address)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  address,
	})
}

// GoEditAddressList 修改地址資料
func (c *AddressController) GoEditAddressList(ctx *gin.Context) {
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	addressId, err := strconv.Atoi(ctx.Query("address_id"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "參數錯誤",
		})
		return
	}
	name := ctx.Query("name")
	phone := ctx.Query("phone")
	address := ctx.Query("address")
	zipcode := ctx.Query("zipcode")
	//將其他地址的預設取消
	model.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	addressModel := model.Address{}
	model.DB.Where("id=?", addressId).Find(&addressModel)
	addressModel.Name = name
	addressModel.Phone = phone
	addressModel.Address = address
	addressModel.Zipcode = zipcode
	addressModel.DefaultAddress = 1 //將修改的地址設為預設
	model.DB.Save(&addressModel)
	allAddress := []model.Address{}
	model.DB.Where("uid=?", user.Id).Order("default_address desc").Find(&allAddress)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  allAddress,
	})
}

// ChangeDefaultAddress 更換預設地址
func (c *AddressController) ChangeDefaultAddress(ctx *gin.Context) {
	user := model.User{}
	model.Cookie.Get(ctx, "userinfo", &user)
	addressId, err := strconv.Atoi(ctx.Query("address_id"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "參數錯誤",
		})
		return
	}
	model.DB.Table("address").Where("uid=?", user.Id).Updates(map[string]interface{}{"default_address": 0})
	model.DB.Table("address").Where("id=?", addressId).Updates(map[string]interface{}{"default_address": 1})
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  "更新預設地址成功",
	})
}
