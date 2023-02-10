package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TokenParse() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenStr := context.Query("token")

		if tokenStr == "" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "没有权限",
			})
		}

		token, claims, err := ParseToken(tokenStr)
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "没有权限",
			})
			context.Abort()
			return
		}
		fmt.Println("claims.UserId")
		fmt.Println(claims.UserId)
		context.Set("my_uid", claims.UserId)
	}
}
