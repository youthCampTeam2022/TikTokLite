package controller

import (
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	service.Response
	VideoList []service.Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {

}
