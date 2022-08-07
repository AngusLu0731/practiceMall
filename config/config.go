package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	AppName              string `yaml:"appName"`
	HttpPort             string `yaml:"httpPort"`
	RunMode              string `yaml:"runMode"`
	AdminPath            string `yaml:"adminPath"`
	ExcludeAuthPath      string `yaml:"excludeAuthPath"`
	SessionOn            bool   `yaml:"sessionOn"`
	SessionGCMaxLifetime int    `yaml:"sessionGCMaxLifetime"`
	sessionName          string `yaml:"sessionName"`
	Domain               string `yaml:"domain"`
	MysqlAdmin           string `yaml:"mysqlAdmin"`
	MysqlPwd             string `yaml:"mysqlPwd"`
	Mysqldb              string `yaml:"mysqldb"`
	MysqlPort            string `yaml:"mysqlPort"`
	MysqlHost            string `yaml:"mysqlHost"`
	ResizeImageSize      []int  `yaml:"resizeImageSize"`
	EnableRedis          string `yaml:"enableRedis"`
	RedisConn            string `yaml:"redisConn"`
	RedisPwd             string `yaml:"redisPwd"`
	RedisTime            int    `yaml:"redisTime"`
	SecureCookie         string `yaml:"secureCookie"`
	CopyRequestBody      bool   `yaml:"copyRequestBody"`
}

var Conf *Config

func init() {
	if f, err := os.Open("config/config.yaml"); err != nil {
		log.Fatalln("Open config.yaml failed : ", err)
	} else {
		err := yaml.NewDecoder(f).Decode(&Conf)
		if err != nil {
			log.Fatalln("Decode yaml failed : ", err)
		}
	}
}
