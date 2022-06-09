package controller

import (
	"TikTokLite/service"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//user:关于用户本身的接口
//登录
//注册
//获取用户信息

//type UserResponse struct {
//	service.Response
//	User service.User `json:"user"`
//}

//UserController 用户注册登录，获取信息控制器
type UserController struct {
	//继承user service服务
	service service.IUserService
}

func NewUserController() *UserController {
	return &UserController{
		service: service.NewUserService(),
	}
}

//返回错误resp
func sendErrResponse(c *gin.Context, resp service.Response) {
	c.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
	})
}

//Register 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var req service.UserLoginOrRegisterRequest
	req.Name, _ = c.GetQuery("username")
	password, _ := c.GetQuery("password")
	if len(password) > 20 {
		resp := service.BuildResponse(errors.New("password is too long,limit <= 20"))
		sendErrResponse(c, resp)
	}
	has := md5.Sum([]byte(password))
	req.Password = fmt.Sprintf("%X", has)
	resp, err := uc.service.UserRegister(&req)
	if err != nil {
		sendErrResponse(c, resp.Response)
		return
	}
	//初始化feed流
	service.UserFeedInit(resp.UserId)
	c.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_id":     resp.UserId,
		"token":       resp.Token,
	})
}

//Login 用户登录
func (uc *UserController) Login(c *gin.Context) {
	var req service.UserLoginOrRegisterRequest
	req.Name, _ = c.GetQuery("username")
	password, _ := c.GetQuery("password")
	has := md5.Sum([]byte(password))
	req.Password = fmt.Sprintf("%X", has)
	resp, err := uc.service.UserLogin(&req)
	if err != nil {
		sendErrResponse(c, resp.Response)
		return
	}
	//初始化feed流
	service.UserFeedInit(resp.UserId)
	c.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_id":     resp.UserId,
		"token":       resp.Token,
	})
}

//UserInfo 获取用户信息
func (uc *UserController) UserInfo(c *gin.Context) {
	var req service.UserInfoRequest
	idStr, _ := c.GetQuery("user_id")
	req.UserId, _ = strconv.ParseInt(idStr, 10, 64)
	req.Token, _ = c.GetQuery("token")
	resp, err := uc.service.UserInfo(&req)
	if err != nil {
		sendErrResponse(c, resp.Response)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user":        resp.User,
	})
}
