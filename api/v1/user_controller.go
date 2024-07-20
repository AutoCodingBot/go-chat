package v1

import (
	"net/http"

	"chat-room/internal/model"
	"chat-room/internal/service"
	"chat-room/internal/utils"
	"chat-room/pkg/common/request"
	"chat-room/pkg/common/response"
	"chat-room/pkg/global/log"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user model.User
	// c.BindJSON(&user)
	c.ShouldBind(&user) //使用报文实体中的数据(json格式)修改struct

	log.Logger.Debug("user", log.Any("user-login-ctl", user))
	err := service.UserService.Login(&user)

	if err != nil {

		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	//generate token
	token, err := utils.GenerateToken(int(user.Id), user.Uuid, user.Username)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
	}
	user.Jwt = token
	c.JSON(http.StatusOK, response.SuccessMsg(&user))
}

func Register(c *gin.Context) {
	var user model.User
	//前端传输的不是form表单,而是json数据,查看USER (struct),可用找到json和model的映射关系!
	c.ShouldBindJSON(&user)
	err := service.UserService.Register(&user)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	//generate token
	token, err := utils.GenerateToken(int(user.Id), user.Uuid, user.Username)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
	}
	c.JSON(200, response.SuccessMsg(token))
}

// 更新头像
func ModifyUserInfo(c *gin.Context) {
	var user model.User
	c.ShouldBindJSON(&user)
	log.Logger.Debug("user", log.Any("user", user))
	if err := service.UserService.ModifyUserInfo(&user); err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, response.SuccessMsg(nil))
}

func GetUserDetails(c *gin.Context) {
	uuid := c.Param("uuid")

	c.JSON(http.StatusOK, response.SuccessMsg(service.UserService.GetUserDetails(uuid)))
}

// 通过用户名获取用户信息
func GetUserOrGroupByName(c *gin.Context) {
	log.Logger.Debug("user", log.Any("user", "In user ctl"))
	name := c.DefaultQuery("name", "")
	if name == "" {
		c.JSON(http.StatusOK, response.FailMsg("搜索词不能为空"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(service.UserService.GetUserOrGroupByName(name)))
}

func GetUserList(c *gin.Context) {
	claims, err := utils.ParseToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uuid := claims.Uuid
	c.JSON(http.StatusOK, response.SuccessMsg(service.UserService.GetUserList(uuid)))
}

func AddFriend(c *gin.Context) {
	var userFriendRequest request.FriendRequest
	c.ShouldBindJSON(&userFriendRequest)

	err := service.UserService.AddFriend(&userFriendRequest)
	if nil != err {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(nil))
}
