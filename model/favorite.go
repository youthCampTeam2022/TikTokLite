package model

import (
	"errors"
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	VideoID int64 `gorm:"index"`
	UserID  int64
}

type VideoRes struct {
	Id            int64   `json:"id,omitempty"`
	Author        UserRes `json:"author"`
	PlayUrl       string  `json:"play_url,omitempty"`
	CoverUrl      string  `json:"cover_url,omitempty"`
	FavoriteCount int64   `json:"favorite_count,omitempty"`
	CommentCount  int64   `json:"comment_count,omitempty"`
	IsFavorite    bool    `json:"is_favorite,omitempty"`
	Title         string  `json:"title,omitempty"`
}

func (f *Favorite) Create() error {
	return DB.Create(&f).Error
}

func (f *Favorite) UniqueInsert() error {
	var FirstRes Favorite
	_ = DB.Model(&Favorite{}).Where("video_id = ? and user_id = ?", f.VideoID, f.UserID).First(&FirstRes).Error
	if FirstRes.ID != 0 {
		return errors.New("repeat favorite")
	}
	return f.Create()
}

func (f *Favorite) Delete() error {
	return DB.Where("user_id=? AND video_id=?", f.UserID, f.VideoID).
		Delete(&Favorite{}).Error
}

func GetFavoriteNum(videoID int64) (count int64) {
	DB.Model(&Favorite{}).Where("video_id = ?", videoID).Count(&count)
	return
}

func (f *Favorite) GetFavoriteList() ([]int64, error) {
	var res []int64
	err := DB.Model(f).Select("video_id").Where("user_id = ?", f.UserID).Scan(&res).Error
	//log.Println(res)
	return res, err
}

func IsFavorite(userId, videoId int64) (bool, error) {
	var count int64
	err := DB.Model(&Favorite{}).Where("video_id = ? and user_id = ?", videoId, userId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func GetFavoriteRes(userID int64) (videos []VideoRes, err error) {
	f := FollowManagerRepository{DB, RedisCache}
	//rows,err := DB.Raw("select favorites.video_id,videos.author_id,videos.play_url,videos.cover_url,videos.title,users.id,users.name" +
	//	"FROM (favorites INNER JOIN videos On favorites.video_id = videos.id) " +
	//	"INNER JOIN users On users.id = videos.author_id" +
	//	"WHERE favorites.deleted_at is null and favorites.user_id = ?",userID).Rows()
	rows, err := DB.Raw("select favorites.video_id,videos.author_id,videos.play_url,videos.cover_url,videos.title "+
		"FROM favorites INNER JOIN videos On favorites.video_id = videos.id "+
		"WHERE favorites.deleted_at is null and favorites.user_id = ?", userID).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var videoRes VideoRes
		err := rows.Scan(&videoRes.Id, &videoRes.Author.Id, &videoRes.PlayUrl,
			&videoRes.CoverUrl, &videoRes.Title)
		if err != nil {
			return nil, err
		}
		videoRes.IsFavorite = true
		videoRes.FavoriteCount = GetFavoriteNum(videoRes.Id)
		videoRes.CommentCount = GetCommentNum(videoRes.Id)
		videoRes.Author = UserRes{
			Id:            videoRes.Author.Id,
			Name:          f.GetName(videoRes.Author.Id),
			FollowCount:   f.RedisFollowCount(videoRes.Author.Id),
			FollowerCount: f.RedisFollowerCount(videoRes.Author.Id),
			IsFollow:      f.RedisIsFollow(userID, videoRes.Author.Id),
		}
		videos = append(videos, videoRes)
	}
	return videos, err
}
