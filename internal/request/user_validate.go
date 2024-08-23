package userreq

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	// NickName        string `json:"nickname" gorm:"nickname"`
	Avatar          string `json:"avatar" gorm:"avatar"`
	CurrentPassword string `json:"currentPassword" gorm:"-"`
	NewPassword     string `json:"newPassword" gorm:"password"`
}

func ValidateUserInfo(c *gin.Context) (UserInfo, error) {
	var userInfo UserInfo
	c.ShouldBindJSON(&userInfo)

	//one of password is empty
	if (userInfo.CurrentPassword == "" && userInfo.NewPassword != "") || (userInfo.CurrentPassword != "" && userInfo.NewPassword == "") {
		error := errors.New("Both password is required")
		return userInfo, error
	}
	//new password
	if userInfo.NewPassword != "" {
		res := validateNewPassword(userInfo.NewPassword)
		if res != true {
			error := errors.New("new password did not suit the case")
			return userInfo, error
		}
	}
	return userInfo, nil
}

func validateNewPassword(newPassword string) bool {
	// 检查长度至少为 6
	if len(newPassword) < 6 {
		return false
	}

	// 检查是否包含小写字母
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(newPassword)
	if !hasLowercase {
		return false
	}

	// 检查是否包含大写字母
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(newPassword)
	if !hasUppercase {
		return false
	}

	// 检查是否包含数字
	hasDigit := regexp.MustCompile(`\d`).MatchString(newPassword)
	if !hasDigit {
		return false
	}

	// 检查是否包含特殊字符
	hasSpecialChar := regexp.MustCompile(`[@#$+-]`).MatchString(newPassword)
	if !hasSpecialChar {
		return false
	}

	return true
}
