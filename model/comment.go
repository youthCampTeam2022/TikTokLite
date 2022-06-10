package model

import (
	"TikTokLite/util"
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
	TotalFavorited int64 `json:"total_favorited,omitempty"`
	FavoriteCount int64 `json:"favorite_count,omitempty"`
	WorkCount int64 `json:"work_count,omitempty"`
}

func GetCommentRes(videoID int64, userID int64) (comments []CommentRes, err error) {
	f := FollowManagerRepository{DB, RedisCache}
	rows, err := DB.Raw("SELECT comments.id,comments.content,comments.created_at,users.id,users.name "+
		"FROM comments INNER JOIN users ON comments.user_id = users.id "+
		"WHERE comments.deleted_at is null and video_id = ? order by comments.id desc", videoID).Rows()
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
		comment.User.TotalFavorited = GetTotalFavoritedRedis(comment.User.Id)
		comment.User.WorkCount = GetTotalWorkCount(comment.User.Id)
		comment.User.FavoriteCount = GetUserFavoriteNum(comment.User.Id)
		comments = append(comments, comment)
	}
	return comments, err
}

func (c *Comment) Create() error {
	err := DB.Create(&c).Error
	if err != nil {
		return err
	}
	IncrCommentRedis(c.VideoID)
	return nil
}

func (c *Comment) Delete() error {
	err := DB.Delete(&c).Error
	if err != nil {
		return err
	}
	DecrCommentRedis(c.VideoID)
	return nil
}

func (c *Comment) DeleteByUser() error {
	DecrCommentRedis(c.VideoID)
	return DB.Where("id=? AND user_id=?", c.ID, c.UserID).Delete(&Comment{}).Error
	//return errors.New("invalid delete")
}

// GetCommentNum 获取评论数
func GetCommentNum(videoID int64) (count int64) {
	DB.Model(&Comment{}).Where("video_id = ?", videoID).Count(&count)
	return
}