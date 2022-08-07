package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"practiceMall/config"
)

var (
	DB      *gorm.DB
	err     error
	sqlConn = config.Conf.MysqlAdmin + ":" + config.Conf.MysqlPwd + "@tcp" + "(" + config.Conf.MysqlHost + ":" + config.Conf.MysqlPort + ")/" + config.Conf.Mysqldb + "?charset=utf8&parseTime=True&loc=Local"
)

func ConnectMySql() {
	DB, err = gorm.Open(mysql.Open(sqlConn), &gorm.Config{})
	if err != nil {
		log.Fatalln("init mysql failed : ", err)
	} else {
		log.Println("init mysql succeed")
	}
}

func CloseMysql() {
	sqlDB, _ := DB.DB()
	sqlDB.Close()
}
