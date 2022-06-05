package controller

import (
	"TikTokLite/middleware"
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	service.Response
	NextTime  int64           `json:"next_time,omitempty"`
	VideoList []service.Video `json:"video_list,omitempty"`
}

// Feed 返回限制时间后发布的n个视频，如果限制时间后没有新的投稿，就按照最新的顺序返回n个视频
func Feed(c *gin.Context) {
	resp := FeedResponse{}
	//返回的信息还得按照用户是否登陆再去判断是否查询is_follow，is_favorite
	_, isLogin := c.GetQuery("token")
	var userId int64 = -1
	if isLogin {
		middleware.ValidDataTokenMiddleWare(c)
		ifUserId, _ := c.Get("user_id")
		userId = ifUserId.(int64)
		//userId, _ = strconv.ParseInt(strUserId.(string), 10, 64)
	}
	latestTime := time.Now().UnixMilli()
	strLatestTime, exist := c.GetQuery("latest_time")
	if exist {
		latestTime, _ = strconv.ParseInt(strLatestTime, 10, 64)
	}
	//videoList, nextTime, err := service.GetFeed(time.UnixMilli(latestTime), userId)
	videoList, nextTime, err := service.GetUserFeed(time.UnixMilli(latestTime), userId)
	if err != nil {
		resp.StatusCode = 1
		resp.StatusMsg = err.Error()
		resp.VideoList = nil
		resp.NextTime = time.Now().UnixMilli()
	} else {
		resp.StatusCode = 0
		resp.StatusMsg = "success"
		resp.VideoList = videoList
		resp.NextTime = nextTime
	}
	c.JSON(http.StatusOK, resp)
}
