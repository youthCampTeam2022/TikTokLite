package controller

import (
	"TikTokLite/model"
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteListResponse struct {
	service.Response
	VideoList []model.VideoRes
}

// FavoriteAction 点赞操作
func FavoriteAction(c *gin.Context) {
	videoIDQuery, _ := c.GetQuery("video_id")
	actionTypeQuery, _ := c.GetQuery("action_type")
	userIDToken, _ := c.Get("user_id")
	userID := userIDToken.(int64)
	videoID, _ := strconv.Atoi(videoIDQuery)
	actionType, _ := strconv.Atoi(actionTypeQuery)

	if actionType == 1 {
		//点赞
		err := service.SetFavorite(int64(videoID), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, service.Response{
				StatusCode: 1,
				StatusMsg:  "err in SetFavorite",
			})
			return
		}

	} else if actionType == 2 {
		//取消点赞
		err := service.CancelFavorite(int64(videoID), userID)
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

// FavoriteList 获取点赞列表
func FavoriteList(c *gin.Context) {
	userIDToken, _ := c.Get("user_id")
	userID := userIDToken.(int64)
	list, err := service.GetFavoriteList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, FavoriteListResponse{
			Response: service.Response{
				StatusCode: 2,
				StatusMsg:  "err in GetFavoriteList",
			},
			VideoList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: service.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		VideoList: list,
	})
}
