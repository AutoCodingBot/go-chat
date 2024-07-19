package utils

import (
	"chat-room/config"
	"chat-room/pkg/global/log"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustClaim struct {
	ID       int    `json:"id"`
	Uuid     string `json:"uuid"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

var secretKey = []byte(config.GetConfig().JwtToken.SecretKey)

func GenerateToken(id int, uuid string, userName string) (string, error) {
	jwtFoo := JwtCustClaim{
		id,
		uuid,
		userName,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GetConfig().JwtToken.TokenExpireTime) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "Token",
		},
	}
	log.Logger.Debug("jwt", log.Any("jwtinfo", config.GetConfig().JwtToken.SecretKey))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtFoo)
	return token.SignedString(secretKey)
	// return res, err
}
func ParseToken(c *gin.Context) (JwtCustClaim, error) {
	tokenStr := c.GetHeader("Authorization")
	parts := strings.Split(tokenStr, " ")
	jwtFoo := JwtCustClaim{}
	token, err := jwt.ParseWithClaims(parts[1], &jwtFoo, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	//token.Valid
	if err != nil && token.Valid {
		err = errors.New(err.Error())
	}
	return jwtFoo, err
}
