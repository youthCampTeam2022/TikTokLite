package model

import (
	"TikTokLite/util"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	gorm.Model
	VideoID int64 `gorm:"index"`
	UserID  int64
	Content string `gorm:"type:varchar(255);"`
}

type CommentRes struct {
	Id         int64   `json:"id,omitempty"`
	User       UserRes `json:"user"`
	Content    string  `json:"content,omitempty"`
	CreateDate string  `json:"create_date,omitempty"`
}

type UserRes struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// GetCommentRes todo: redis能正常用的话这部分可以不用联查了
func GetCommentRes(videoID int64, userID int64) (comments []CommentRes, err error) {
	f := FollowManagerRepository{DB, RedisCache}
	rows, err := DB.Raw("SELECT comments.id,comments.content,comments.created_at,users.id,users.name "+
		"FROM comments INNER JOIN users ON comments.user_id = users.id "+
		"WHERE comments.deleted_at is null and video_id = ?", videoID).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var comment CommentRes
		var createDate time.Time
		err := rows.Scan(&comment.Id, &comment.Content, &createDate, &comment.User.Id, &comment.User.Name)
		if err != nil {
			return nil, err
		}
		comment.CreateDate = util.Time2String(createDate)
		comment.User.FollowerCount = f.RedisFollowCount(comment.User.Id)
		comment.User.FollowCount = f.RedisFollowerCount(comment.User.Id)
		comment.User.IsFollow = f.RedisIsFollow(userID, comment.User.Id)
		comments = append(comments, comment)
	}
	return comments, err
}

func (c *Comment) Create() error {
	return DB.Create(&c).Error
}

func (c *Comment) Delete() error {
	return DB.Delete(&c).Error
}

func (c *Comment) DeleteByUser() error {
	var uid Comment
	DB.Model(&c).First(&uid)
	if uid.UserID == c.UserID {
		return DB.Delete(&c).Error
	}
	return errors.New("invalid delete")
}

func GetCommentNum(videoID int64) (count int64) {
	DB.Model(&Comment{}).Where("video_id = ?", videoID).Count(&count)
	return
}

func GetCommentsByVideo(videoID int64) (comments []Comment, err error) {
	err = DB.Model(&Comment{}).Where("video_id = ?", videoID).Find(&comments).Order("created_at DESC").Error
	return
}
