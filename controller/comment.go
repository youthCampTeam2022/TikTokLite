package controller

import (
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	service.Response
	CommentList []service.Comment `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {

}

func CommentList(c *gin.Context) {

}
