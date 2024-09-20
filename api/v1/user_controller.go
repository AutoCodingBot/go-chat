package v1

import (
	"net/http"

	"chat-room/internal/model"
	userreq "chat-room/internal/request"
	"chat-room/internal/service"
	"chat-room/internal/utils"
	"chat-room/pkg/common/request"
	"chat-room/pkg/common/response"
	"chat-room/pkg/global/log"
	"chat-room/pkg/global/state"

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
// func ModifyUserInfo(c *gin.Context) {
// 	var user model.User
// 	c.ShouldBindJSON(&user)
// 	log.Logger.Debug("user", log.Any("user", user))
// 	if err := service.UserService.ModifyUserInfo(&user); err != nil {
// 		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusNoContent, response.SuccessMsg(nil))
// }

func GetUserDetails(c *gin.Context) {
	claims, err := utils.ParseToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(service.UserService.GetUserDetails(claims.Uuid)))
}

func GetUserOrGroupByName(c *gin.Context) {
	// log.Logger.Debug("user", log.Any("user", "In user ctl"))
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
	friendsSlice := service.UserService.GetUserList(claims)
	for index, val := range friendsSlice {
		friendsSlice[index].OnlineStatus = state.UserOnlineStatus(val.Uuid)
	}
	c.JSON(http.StatusOK, response.SuccessMsg(friendsSlice))
}

func AddFriend(c *gin.Context) {
	claims, err := utils.ParseToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var userFriendRequest request.FriendRequest
	c.ShouldBindJSON(&userFriendRequest)

	data, err := service.UserService.AddFriend(claims, &userFriendRequest)
	if nil != err {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(data))
}

func UpdateUserProfile(c *gin.Context) {
	// log.Logger.Debug("user", log.Any("user", userInfo))
	// log.Logger.Debug("user", log.Any("err", err))
	claims, err := utils.ParseToken(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	//validate
	userInfo, err := userreq.ValidateUserInfo(c)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	//update
	err = service.UserService.UpdateUserProfile(claims, &userInfo)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessMsg(userInfo))
}

func UserOnlineStatus(c *gin.Context) {
	log.Logger.Debug("debug", log.Any("debug", "hit me"))

	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing UUID parameter"})
		return
	}
	useOnlineStatus := state.UserOnlineStatus(uuid)
	c.JSON(http.StatusOK, response.SuccessMsg(useOnlineStatus))
}
