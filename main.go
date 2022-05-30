package main

import (
	"TikTokLite/model"
	"TikTokLite/router"
	"TikTokLite/setting"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := setting.Init("./config/config.yaml"); err != nil {
		fmt.Printf("init setting failed, err: %v \n", err)
		return
	}
	fmt.Printf("\n%v\n", setting.Conf)
	model.Init()
	r := gin.Default()
	router.RouterInit(r)
	r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
}
