package controller

import (
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
)

//user:关于用户本身的接口
//登录
//注册
//获取用户信息

type UserLoginResponse struct {
	service.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	service.Response
	User service.User `json:"user"`
}


func Register(c *gin.Context) {

}

func Login(c *gin.Context) {

}

func UserInfo(c *gin.Context) {

}