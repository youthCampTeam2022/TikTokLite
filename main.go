package main

import (
	"TikTokLite/model"
	"TikTokLite/router"
	"github.com/gin-gonic/gin"
)

func main()  {
	model.MysqlInit()
	r := gin.Default()
	router.RouterInit(r)
	r.Run()
}
