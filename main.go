package main

import (
	"TikTokLite/model"
	"TikTokLite/router"
	"TikTokLite/setting"
	"TikTokLite/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
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
	//定时更新hotfeed和推送
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			model.BuildHotFeed()
			model.CheckAliveUserAndPushHotFeed()
		}
	}()
	r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
}
