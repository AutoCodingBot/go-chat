package v1

import (
	"net/http"
	"strconv"

	"chat-room/internal/service"
	"chat-room/internal/utils"
	"chat-room/pkg/common/response"
	"chat-room/pkg/global/state"

	"github.com/gin-gonic/gin"
)

type messageResponseConcat struct {
	MessageResponse *[]response.MessageResponse `json:"messageList"`
	OnlineStatus    bool                        `json:"onlineStatus"`
	// LatestMsg       string                      `json:"latestMsg"`
}

// 获取消息列表
func GetMessage(c *gin.Context) {
	/*
		Uuid: 28353ed6-5966-4804-9c52-9b00abd4401e
		FriendUsername: sam
		MessageType: 1
	*/
	claims, err := utils.ParseToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	currentUserId := claims.ID
	currentUserName := claims.UserName
	messageTypeStr := c.Query("MessageType")
	messageType, err := strconv.Atoi(messageTypeStr)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
	}
	friendUsername := c.Query("FriendUsername")
	freindUuid := c.Query("Uuid")
	messages, err := service.MessageService.GetMessages(currentUserId, currentUserName, friendUsername, messageType)

	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}
	//用户聊天消息
	if messageTypeStr == "1" {
		res := messageResponseConcat{}
		res.MessageResponse = &messages //消息记录
		//更新在线详情
		res.OnlineStatus = state.UserOnlineStatus(freindUuid)
		c.JSON(http.StatusOK, response.SuccessMsg(res))
	} else {
		//群消息
		c.JSON(http.StatusOK, response.SuccessMsg(messages))

	}

}
