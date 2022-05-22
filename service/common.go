package service

import (
	"TikTokLite/model"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

const (
	SecretKey = "0SFWF023423dhwq"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	//videos表
	Id int64 `json:"id,omitempty"`
	//users表
	Author User `json:"author"`
	//videos
	PlayUrl  string `json:"play_url,omitempty"`
	CoverUrl string `json:"cover_url,omitempty"`
	//favorites
	FavoriteCount int64 `json:"favorite_count,omitempty"`
	//comments
	CommentCount int64 `json:"comment_count,omitempty"`
	//favorites
	IsFavorite bool `json:"is_favorite,omitempty"`
	//videos
	Title string `json:"title,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

//BuildUserList 通用接口，传入user_id集合和follow仓库接口，返回[]User（Id，Name，FollowCount，FollowerCount，IsFollow）
func BuildUserList(userID int64, userIDList []int64, m model.IFollowRepository) []User {
	Users := make([]User, len(userIDList))
	for i := 0; i < len(userIDList); i++ {
		Users[i] = BuildUser(userID, userIDList[i], m)
	}
	return Users
}

/*BuildUser 返回User（Id，Name，FollowCount，FollowerCount，IsFollow）
userID是当前用户的id，toUserID是要查询的ID，m是follow仓库的接口*/
func BuildUser(userID, toUserID int64, m model.IFollowRepository) User {
	var user User
	user.Id = toUserID
	user.Name = m.GetName(toUserID)
	user.IsFollow = m.RedisIsFollow(userID, toUserID)
	user.FollowCount = m.RedisFollowCount(toUserID)
	user.FollowerCount = m.RedisFollowerCount(toUserID)
	return user
}

//BuildResponse 根据err返回对应的Response
func BuildResponse(err error) Response {
	resp := Response{}
	if err != nil {
		//这里暂时还没协商出errno code，所以错误先默认为
		log.Println(err)
		resp.StatusCode = -1
		resp.StatusMsg = "fail"
	} else {
		resp.StatusCode = 0
		resp.StatusMsg = "success"
	}
	return resp
}

//Claims token claims
type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.StandardClaims
}

//GetToken 获取对应user的token
func GetToken(u *model.User) (tokenString string, rep Response, err error) {
	var claims Claims
	claims.Username = u.Name
	claims.UserID = int64(u.ID)
	//token过期时间
	claims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	//token创建时间
	claims.IssuedAt = time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密生成token
	tokenString, err = token.SignedString([]byte(SecretKey))
	rep = BuildResponse(err)
	return
}
