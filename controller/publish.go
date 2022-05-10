package controller

import (
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	service.Response
	VideoList []service.Video `json:"video_list"`
}

func Publish(c *gin.Context) {

}

func PublishList(c *gin.Context) {

}