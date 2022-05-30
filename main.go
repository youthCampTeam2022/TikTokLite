package main

import (
	"TikTokLite/model"
	"TikTokLite/router"
	"TikTokLite/util"
	"github.com/gin-gonic/gin"
)

func main() {
	model.Init()
	r := gin.Default()
	router.RouterInit(r)
	util.FilterInit()
	r.Run(":8081")
}
