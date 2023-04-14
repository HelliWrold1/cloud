package middleware

import (
	"github.com/HelliWrold1/cloud/internal/crypt"
	"github.com/HelliWrold1/cloud/internal/ecode"
	"github.com/gin-gonic/gin"
	"time"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		code = 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claims, err := crypt.ParseToken(token)
			if err != nil {
				code = ecode.AccessDenied.Code()
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = ecode.DeadlineExceeded.Code()
			}
		}
		if code != 200 {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    ecode.AccessDenied.Msg(),
				"data":   data,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
