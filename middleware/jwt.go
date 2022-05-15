package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SecretKey = "0SFWF0234-432WFSB-EFRHG34234-432WFDEN-dserhwe-342423bjfds-342jfdxj320r324-3bfjsdbj"
)

//ValidDataTokenMiddleWare token鉴权中间件
func ValidDataTokenMiddleWare(c *gin.Context) {
	tokenString, _ := c.GetQuery("token")
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status_code": -1,
			"status_msg":  "unauthorized access",
		})
		c.Abort()
		return
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_name", claims["user_name"])
			c.Set("user_id", claims["user_id"])
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
