package main

import (
	"TikTokLite/model"
	"TikTokLite/router"
	"github.com/gin-gonic/gin"
)

func main() {
	model.Init()
	r := gin.Default()
	router.RouterInit(r)
	r.Run(":8081")
}
