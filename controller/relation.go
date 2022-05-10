package controller

import (
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	service.Response
	UserList []service.User `json:"user_list"`
}

func RelationAction(c *gin.Context) {

}

func FollowList(c *gin.Context) {

}

func FollowerList(c *gin.Context) {

}