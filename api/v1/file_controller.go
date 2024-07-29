package v1

import (
	"io/ioutil"
	"net/http"
	"strings"

	"chat-room/config"
	"chat-room/internal/utils"
	"chat-room/pkg/common/response"
	"chat-room/pkg/global/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 前端通过文件名称获取文件流，显示文件
func GetFile(c *gin.Context) {
	fileName := c.Param("fileName")
	data, _ := ioutil.ReadFile(config.GetConfig().StaticPath.FilePath + fileName)
	c.Writer.Header().Set("Cache-Controler", "public,max-age=86400")
	c.Writer.Write(data)
}

// 上传头像等文件
func SaveFile(c *gin.Context) {
	claims, err := utils.ParseToken(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	namePreffix := uuid.New().String()

	// userUuid := claims.Uuid
	// objectType := c.PostForm("objectType")

	file, _ := c.FormFile("file")
	fileName := file.Filename
	index := strings.LastIndex(fileName, ".")
	suffix := fileName[index:]

	newFileName := namePreffix + suffix

	log.Logger.Info("file", log.Any("file avatar addr", config.GetConfig().StaticPath.FilePath+newFileName))
	log.Logger.Info("file", log.Any("user change avatar,user:", claims.UserName))

	err = c.SaveUploadedFile(file, config.GetConfig().StaticPath.FilePath+newFileName)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(newFileName))
	// return
	// err = service.UserService.ModifyUserAvatar(newFileName, claims.ID, objectType)
	// if err != nil {
	// 	c.JSON(http.StatusOK, response.FailMsg(err.Error()))
	// }
	// c.JSON(http.StatusOK, response.SuccessMsg(newFileName))
}
