package main

import (
	"TikTokLite/model"
	"TikTokLite/router"
	"TikTokLite/service"
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
	service.UpdateUnLoginFeed()
	//定时更新hotfeed和推送
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			model.BuildHotFeed()
			model.CheckAliveUserAndPushHotFeed()
			service.UpdateUnLoginFeed()
		}
	}()
	r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
}
