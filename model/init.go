package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	DB         *gorm.DB
	RedisCache *Cache
)

//注意不要明文放github上
const (
	user     = ""
	password = ""
)

func MysqlInit() {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/tiktoklite?charset=utf8mb4&parseTime=True&loc=Local", user, password)
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
	RedisCache = NewRedisCache(0, "127.0.0.1:6379", FOREVER)
}
func Init() {
	MysqlInit()
	RedisInit()
}
