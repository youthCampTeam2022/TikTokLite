package model

import (
	"TikTokLite/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	DB         *gorm.DB
	RedisCache *Cache
)

func MysqlInit() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", setting.Conf.MysqlConfig.User, setting.Conf.MysqlConfig.Password, setting.Conf.MysqlConfig.Host, setting.Conf.MysqlConfig.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("err in MysqlInit:", err)
		return
	}
	DB = db
	//fmt.Println(DB)
	_ = DB.AutoMigrate(&Video{})
	_ = DB.AutoMigrate(&User{})
	_ = DB.AutoMigrate(&Follow{})
	_ = DB.AutoMigrate(&Comment{})
	_ = DB.AutoMigrate(&Favorite{})

}
func RedisInit() {
	RedisCache = NewRedisCache(setting.Conf.DB, setting.Conf.RedisConfig.Host, FOREVER)
}
func Init() {
	MysqlInit()
	RedisInit()
}
