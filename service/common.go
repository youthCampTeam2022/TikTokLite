package service

import (
	"TikTokLite/model"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var SecretKey = []byte("djwlqjsk-dwqjdk2k3u-vmsdmw-342f-ewrk-nk23u2i4")

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
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

//GetToken 获取对应user的token
func GetToken(u *model.User) (tokenString string, rep Response, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	var claims = make(jwt.MapClaims)
	claims["user_name"] = u.Name
	claims["user_id"] = u.ID
	//token过期时间
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(24)).Unix()
	//token创建时间
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	//加密生成token
	tokenString, err = token.SignedString([]byte(SecretKey))
	rep = BuildResponse(err)
	return
}
