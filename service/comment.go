package service

import (
	"TikTokLite/model"
	"TikTokLite/util"
	"gorm.io/gorm"
)

func GetComment(videoID int64) ([]Comment, error) {
	comments, err := model.GetCommentsByVideo(videoID)
	if err != nil {
		return nil, err
	}
	commentResult := make([]Comment, len(comments))
	for i, comment := range comments {
		u := new(model.User)
		err := model.NewUserManagerRepository().GetById(u, uint(comment.UserID))
		if err != nil {
			return nil, err
		}
		commentResult[i] = Comment{
			Id: int64(comment.ID),
			User: User{
				Id:   int64(u.ID),
				Name: u.Name,
				//todo: 缺接口
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      false,
			},
			Content:    comment.Content,
			CreateDate: util.Time2String(comment.CreatedAt),
		}
	}
	return commentResult, nil
}

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
