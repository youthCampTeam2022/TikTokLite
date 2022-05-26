package middleware

import (
	"TikTokLite/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//ValidDataTokenMiddleWare token鉴权中间件
func ValidDataTokenMiddleWare(c *gin.Context) {
	tokenString, exist := c.GetQuery("token")
	//投稿接口上的token是放在表单里面的
	if !exist {
		tokenString = c.PostForm("token")
	}
	token, err := jwt.ParseWithClaims(tokenString, &service.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.SecretKey), nil
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code": -1,
			"status_msg":  "unauthorized access",
		})
		c.Abort()
		return
	} else {
		if claims, ok := token.Claims.(*service.Claims); ok && token.Valid {
			c.Set("user_name", claims.Username)
			c.Set("user_id", claims.UserID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": -1,
				"status_msg":  "token is not valid",
			})
			c.Abort()
			return
		}
	}
}
