package middleware

import (
	"chat-room/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(config.GetConfig().JwtToken.SecretKey)

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
	//token有效?
	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		// 验证Signing Method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
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
