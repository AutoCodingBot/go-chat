package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	//获取报文header中的Authorization
	tokenStr := c.GetHeader("Authorization")
	//拆分@param:tokenStr
	parts := strings.Split(tokenStr, " ")
	//判断Authorization 是否合规
	if parts[0] != "Bearer" || len(parts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
		c.Abort()
		return
	}
	c.Next()

	/*
		//判断Authorization 是否有效
		jwtPayload, err := utils.ParseToken(parts[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		}
		// log.Logger.Debug("telegram", log.Any("token-data", data))

		//返回结构体(id,username...)
		c.Set("jwtPayload", jwtPayload) // 设置键为"jwtPayload"的上下文值
		c.Next()                        // 调用下一个处理器
	*/
}
