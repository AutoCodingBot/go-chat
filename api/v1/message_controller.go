package v1

import (
	"net/http"

	"chat-room/internal/service"
	"chat-room/pkg/common/request"
	"chat-room/pkg/common/response"
	"chat-room/pkg/global/log"

	"github.com/gin-gonic/gin"
)

// 获取消息列表
func GetMessage(c *gin.Context) {
	/*
		Uuid: 28353ed6-5966-4804-9c52-9b00abd4401e   个人聊天:uuid=>用户uuid;群:uuid=>group uuid
		FriendUsername: eric
		MessageType: 1
	*/
	log.Logger.Info(c.Query("uuid"))
	var messageRequest request.MessageRequest
	err := c.BindQuery(&messageRequest)
	if nil != err {
		log.Logger.Error("bindQueryError", log.Any("bindQueryError", err))
	}
	log.Logger.Info("messageRequest params: ", log.Any("messageRequest", messageRequest))

	messages, err := service.MessageService.GetMessages(messageRequest)
	if err != nil {
		c.JSON(http.StatusOK, response.FailMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessMsg(messages))
}
