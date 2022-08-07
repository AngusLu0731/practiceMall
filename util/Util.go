package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gomarkdown/markdown"
	_ "github.com/jinzhu/gorm"
	"io/ioutil"
	"math/rand"
	"path"
	"practiceMall/config"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimestampToDate 時間轉日期
func TimestampToDate(ts int) string {
	t := time.Unix(int64(ts), 0)
	return t.Format("2006-04-02 15-04-05")
}

// GetUnix 獲取Timestamp
func GetUnix() int64 {
	return time.Now().Unix()
}

// GetUnixNano 獲取Timestamp的Nano時間
func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

// GetDate 獲取當前日期和時間
func GetDate() string {
	tmp := "2006-04-02 15-04-05"
	return time.Now().Format(tmp)
}

// FormatDay 獲取當前日期
func FormatDay() string {
	tmp := "20060402"
	return time.Now().Format(tmp)
}

// Md5 加密
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// VerifyEmail 驗證Email
func VerifyEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// GenerateOrderId 生成訂單號
func GenerateOrderId() string {
	tmp := "200604021504"
	return time.Now().Format(tmp)
}

// ResizeImage 裁切圖片
func ResizeImage(filename string) {
	extName := path.Ext(filename)
	for i := 0; i < len(config.Conf.ResizeImageSize); i++ {
		width := config.Conf.ResizeImageSize[i]
		savePath := filename + "_" + string(width) + "x" + string(width) + extName
		_ = savePath
	}
}

// FormatAttr 格式化次級標題
func FormatAttr(str string) string {
	md := []byte(str)
	htmlByte := markdown.ToHTML(md, nil, nil)
	return string(htmlByte)
}

// Mul 乘法
func Mul(f float64, num int) float64 {
	return f * float64(num)
}

// GetRandomNum 生產隨機4位數
func GetRandomNum() string {
	var str string
	for i := 0; i < 4; i++ {
		current := rand.Intn(10)
		str += strconv.Itoa(current)
	}
	return str
}

// SendMsg 發送驗證碼
func SendMsg(str string) {
	// 驗證碼需申請，目前先固定值
	ioutil.WriteFile("test_send.txt", []byte(str), 06666)
}

func FormatImg(picName string) string {
	flag := strings.Contains(picName, "/static")
	if flag {
		return picName
	}
	return "/" + picName
}

func SubStr(str string, start int, length int) string {
	if start+length <= len(str) {
		return str[start : start+length]
	}
	return str
}

func StringInSlice(str string, list []string) bool {
	for _, a := range list {
		if a == str {
			return true
		}
	}
	return false
}
