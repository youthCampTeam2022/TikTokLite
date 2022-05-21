package controller

import (
	"TikTokLite/model"
	"TikTokLite/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	service.Response
	CommentList []model.CommentRes `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {
	//token := c.Query("token")
	userIDQuery,_ := c.GetQuery("user_id")
	videoIDQuery, _ := c.GetQuery("video_id")
	actionTypeQuery,_ := c.GetQuery("action_type")
	commentTextQuery,_ := c.GetQuery("comment_text")
	commentIDQuery,_ := c.GetQuery("comment_id")

	videoID, _ := strconv.Atoi(videoIDQuery)
	userID,_ := strconv.Atoi(userIDQuery)
	commentID, _ := strconv.Atoi(commentIDQuery)
	actionType,_ := strconv.Atoi(actionTypeQuery)

	if actionType==1 {
		comment,ok := service.CommentFilter(commentTextQuery)
		if !ok{
			c.JSON(http.StatusInternalServerError,service.Response{
				StatusCode: 4,
				StatusMsg:  "too many sensitive words",
			})
			return
		}
		err := service.CreateComment(int64(videoID), int64(userID),comment)
		if err != nil {
			c.JSON(http.StatusInternalServerError,service.Response{
				StatusCode: 5,
				StatusMsg:  "err in CreateComment",
			})
			return
		}
	}else if actionType == 2{
		err := service.DeleteComment(int64(userID), int64(commentID))
		if err != nil {
			c.JSON(http.StatusInternalServerError,service.Response{
				StatusCode: 6,
				StatusMsg:  "err in CreateComment",
			})
			return
		}
	}else {
		c.JSON(http.StatusInternalServerError,service.Response{
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

func CommentList(c *gin.Context) {
	//token := c.Query("token")
	//userIDQuery,_ := c.GetQuery("user_id")
	videoIDQuery,_  := c.GetQuery("video_id")
	//过验证
	videoID, _ := strconv.Atoi(videoIDQuery)
	//userID, _ := strconv.Atoi(userIDQuery)
	//comments,err := service.GetComment(int64(videoID))
	comments,err := service.GetCommentByJoin(int64(videoID))
	if err != nil{
		c.JSON(http.StatusInternalServerError,CommentListResponse{
			Response:    service.Response{StatusCode: 1, StatusMsg: "err in get comment"},
			CommentList: []model.CommentRes{},
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    service.Response{StatusCode: 0},
		CommentList: comments,
	})
}
