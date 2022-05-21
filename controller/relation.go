package controller

import (
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RelationController struct {
	service service.IFollowService
}

func NewRelationController() *RelationController {
	return &RelationController{
		service: service.NewFollowService(),
	}
}

type UserListResponse struct {
	service.Response
	UserList []service.User `json:"user_list"`
}

func (rc *RelationController) RelationAction(c *gin.Context) {
	var req service.RelationActionRequest
	req.Token, _ = c.GetQuery("token")
	toUserIDS, _ := c.GetQuery("to_user_id")
	actionTpS, _ := c.GetQuery("action_type")
	actionTp, _ := strconv.ParseInt(actionTpS, 10, 32)
	req.ActionType = int32(actionTp)
	req.ToUserID, _ = strconv.ParseInt(toUserIDS, 10, 64)
	userID, _ := c.Get("user_id")
	req.UserID = userID.(int64)
	resp, err := rc.service.RedisRelationAction(&req)
	if err != nil {
		sendErrResponse(c, resp.Response)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
	})

}

//BindFollowListRequest 从url中读取参数
func BindFollowListRequest(c *gin.Context) service.RelationFollowListRequest {
	var req service.RelationFollowListRequest
	req.Token, _ = c.GetQuery("token")
	userIDS, _ := c.GetQuery("user_id")
	req.UserID, _ = strconv.ParseInt(userIDS, 10, 64)
	return req
}

func (rc *RelationController) FollowList(c *gin.Context) {
	req := BindFollowListRequest(c)
	resp, err := rc.service.RedisFollowList(req.UserID)
	if err != nil {
		sendErrResponse(c, resp.Response)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"user_list":   resp.Users,
	})
}

func (rc *RelationController) FollowerList(c *gin.Context) {
	req := BindFollowListRequest(c)
	resp, err := rc.service.RedisFollowerList(req.UserID)
	if err != nil {
		sendErrResponse(c, resp.Response)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"user_list":   resp.Users,
	})
}
