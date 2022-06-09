package controller

import (
	"TikTokLite/model"
	"TikTokLite/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	service.Response
	VideoList []service.Video `json:"video_list"`
}

func ResponseError(c *gin.Context, code int32, err error) {
	c.JSON(http.StatusOK, service.Response{
		StatusCode: code,
		StatusMsg:  err.Error(),
	})
}

func Publish(c *gin.Context) {
	userId, _ := c.Get("user_id")
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		ResponseError(c, 1, err)
		return
	}
	video := model.Video{
		AuthorId: userId.(int64),
		Title:    title,
	}
	//截取封面，上传视频和封面并返回外链
	err = service.PublishVideo(data, userId.(int64), video, c)
	if err != nil {
		log.Println("publish failed，err：", err)
		ResponseError(c, 1, err)
	}
	c.JSON(http.StatusOK, service.Response{
		StatusCode: 0,
		StatusMsg:  fmt.Sprintf("userID:%d,title:%s, uploaded successfully", userId, title),
	})
}

func PublishList(c *gin.Context) {
	//token 里面解析出来的为发出请求的用户id
	userId, _ := c.Get("user_id")
	//query中的user_id才是需要被查询的用户id
	strToUserId, _ := c.GetQuery("user_id")
	toUserId, _ := strconv.ParseInt(strToUserId, 10, 64)
	sResp := service.Response{}
	resp := VideoListResponse{}
	videoList, err := service.GetVideoList(userId.(int64), toUserId)
	if err != nil {
		sResp.StatusCode = 1
		sResp.StatusMsg = "getVideoList failed"
		resp.VideoList = nil
	} else {
		sResp.StatusCode = 0
		sResp.StatusMsg = "success"
		resp.VideoList = videoList
	}
	resp.Response = sResp
	c.JSON(http.StatusOK, resp)
}
