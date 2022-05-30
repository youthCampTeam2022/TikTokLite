package main

import (
	"TikTokLite/model"
	"TikTokLite/router"
	"TikTokLite/setting"
	"TikTokLite/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := setting.Init("./config/config.yaml"); err != nil {
		fmt.Printf("init setting failed, err: %v \n", err)
		return
	}
	model.Init()
	r := gin.Default()
	router.RouterInit(r)
	util.FilterInit()
	r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
}
