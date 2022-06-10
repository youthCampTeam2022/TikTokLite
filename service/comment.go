package service

import (
	"TikTokLite/model"
	"TikTokLite/util"
	"gorm.io/gorm"
)

// GetCommentByJoin 改用联查的版本
func GetCommentByJoin(videoID int64, userID int64) ([]model.CommentRes, error) {
	return model.GetCommentRes(videoID, userID)
}

func CreateComment(videoID, userID int64, text string) (model.Comment, error) {
	c := model.Comment{
		VideoID: videoID,
		UserID:  userID,
		Content: text,
	}
	return c, c.Create()
}

func DeleteComment(userID, commentID, videoID int64) error {
	c := model.Comment{
		Model: gorm.Model{
			ID: uint(commentID),
		},
		UserID:  userID,
		VideoID: videoID,
	}
	return c.DeleteByUser()
}

// CommentFilter 评论过滤器，过滤敏感词
func CommentFilter(commentMsg string) (string, bool) {
	return util.Filtration(commentMsg)
}
