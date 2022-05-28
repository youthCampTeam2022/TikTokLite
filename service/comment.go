package service

import (
	"TikTokLite/model"
	"gorm.io/gorm"
)

// GetCommentByJoin 改用联查的版本，原来的太蠢了
func GetCommentByJoin(videoID int64, userID int64) ([]model.CommentRes, error) {
	return model.GetCommentRes(videoID, userID)
}

func CreateComment(videoID, userID int64, text string) error {
	c := model.Comment{
		VideoID: videoID,
		UserID:  userID,
		Content: text,
	}
	return c.Create()
}

func CommentFilter(commentMsg string) (string, bool) {
	return commentMsg, true
}

func DeleteComment(userID, commentID int64) error {
	c := model.Comment{
		Model: gorm.Model{
			ID: uint(commentID),
		},
		UserID: userID,
	}
	return c.DeleteByUser()
}
