package controller

import (
	"TikTokLite/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteListResponse struct {
	service.Response
	videoList []service.Video
}

func FavoriteAction(c *gin.Context) {
	userIDQuery, _ := c.GetQuery("user_id")
	videoIDQuery, _ := c.GetQuery("video_id")
	actionTypeQuery, _ := c.GetQuery("action_type")

	videoID, _ := strconv.Atoi(videoIDQuery)
	userID, _ := strconv.Atoi(userIDQuery)
	actionType, _ := strconv.Atoi(actionTypeQuery)

	if actionType == 1 {
		//点赞
		fmt.Println(videoID, userID)
		err := service.SetFavorite(int64(videoID), int64(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, service.Response{
				StatusCode: 1,
				StatusMsg:  "err in SetFavorite",
			})
			return
		}

	} else if actionType == 2 {
		//取消点赞
		err := service.CancelFavorite(int64(videoID), int64(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, service.Response{
				StatusCode: 2,
				StatusMsg:  "err in CancelFavorite",
			})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, service.Response{
			StatusCode: 3,
			StatusMsg:  "invalid action_type",
		})
		return
	}
	c.JSON(http.StatusOK, service.Response{
		StatusCode: 0,
		StatusMsg:  "ok",
	})
}

func FavoriteList(c *gin.Context) {
	//token和user
	userIDQuery, _ := c.GetQuery("user_id")
	userID, _ := strconv.Atoi(userIDQuery)
	list, err := service.GetFavoriteList(int64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, FavoriteListResponse{
			Response: service.Response{
				StatusCode: 2,
				StatusMsg:  "err in GetFavoriteList",
			},
			videoList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: service.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		videoList: list,
	})
}
