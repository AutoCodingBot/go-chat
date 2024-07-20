package service

import (
	"chat-room/internal/dao/pool"
	"chat-room/internal/utils"
	"chat-room/pkg/common/constant"
	"chat-room/pkg/common/response"
	"chat-room/pkg/errors"
	"chat-room/pkg/global/log"
	"chat-room/pkg/protocol"

	"chat-room/internal/model"

	"gorm.io/gorm"
)

const NULL_ID int32 = 0

type messageService struct {
}

var MessageService = new(messageService)

func (m *messageService) GetMessages(currentUserId int, currentUserName string, friendUsername string, messageType int) ([]response.MessageResponse, error) {
	db := pool.GetDB()

	migrate := &model.UserMessage{}
	pool.GetDB().AutoMigrate(&migrate)

	if messageType == constant.MESSAGE_TYPE_USER {
		//Friend Info
		var friendInfo *model.User
		db.Select("id").First(&friendInfo, "username = ?", friendUsername)

		if NULL_ID == friendInfo.Id {
			return nil, errors.New("U did not have this friend!")
		}

		var messages []response.MessageResponse

		//generate conversation_id
		conversationId := utils.GenerateConversationId(currentUserName, friendUsername)
		//分步查询
		//一:获取主表内容
		// var userMessageModel model.UserMessage
		// db.Debug().Select("*").Model(userMessageModel).Where("conversation_id = ?", conversationId).Find(&messages)
		//补全avatar,senderUserName,receiverUserName
		//二:获取用户信息
		//三:遍历主表,通过用户信息补全
		// log.Logger.Info("Dot messages", log.Any("Dot messages", messages))

		db.Table("user_messages AS m").
			Select("m.id, m.from_user_id, m.to_user_id, m.content, m.content_type, m.url, m.created_at, u.username AS from_username, u.avatar, to_user.username AS to_username").
			Joins("LEFT JOIN users AS u ON m.from_user_id = u.id").
			Joins("LEFT JOIN users AS to_user ON m.to_user_id = to_user.id").
			Where("m.conversation_id = ?", conversationId).
			Scan(&messages)
		return messages, nil
	}

	if messageType == constant.MESSAGE_TYPE_GROUP {
		messages, err := fetchGroupMessage(db, currentUserId, friendUsername)
		if err != nil {
			return nil, err
		}

		return messages, nil
	}

	return nil, errors.New("不支持查询类型")
}

func fetchGroupMessage(db *gorm.DB, currentUserId int, friendUsername string) ([]response.MessageResponse, error) {
	var group model.Group
	db.First(&group, "name = ?", friendUsername)
	if group.ID <= 0 {
		return nil, errors.New("群组不存在")
	}
	//  Is current User in group
	var groupMember model.GroupMember
	var counts int64
	db.Model(groupMember).Where("user_id=?", currentUserId).Count(&counts)
	if counts == 0 {
		return nil, errors.New("U are not in this group")
	}
	var messages []response.MessageResponse

	db.Raw("SELECT gm.*,u.username,u.avatar from group_messages  `gm` left join users `u` on u.id = gm.from_user_id where group_id = ?",
		group.ID).Scan(&messages)
	// log.Logger.Info("Group messages", log.Any("none", messages))

	return messages, nil
}

func (m *messageService) SaveMessage(message protocol.Message) {
	db := pool.GetDB()
	var fromUser model.User
	db.Select("id,username").Find(&fromUser, "uuid = ?", message.From)
	if NULL_ID == fromUser.Id {
		log.Logger.Error("SaveMessage not find from user", log.Any("SaveMessage not find from user", fromUser.Id))
		return
	}

	var toUserId int32 = 0

	if message.MessageType == constant.MESSAGE_TYPE_USER {
		var toUser model.User

		//reveiver user exists?
		db.Find(&toUser, "uuid = ?", message.To)
		if NULL_ID == toUser.Id {
			return
		}
		toUserId = toUser.Id
		//generate ConversationId
		conversationId := utils.GenerateConversationId(toUser.Username, fromUser.Username)
		saveMessage := model.UserMessage{
			FromUserId:     fromUser.Id,
			ToUserId:       toUserId,
			Content:        message.Content,
			ContentType:    int16(message.ContentType),
			Url:            message.Url,
			ConversationId: conversationId,
		}
		db.Save(&saveMessage)
	}

	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		var group model.Group
		db.Find(&group, "uuid = ?", message.To)
		if NULL_ID == group.ID {
			return
		}
		// ID := group.ID
		saveMessage := model.GroupMessage{
			GroupId:     group.ID,
			FromUserId:  fromUser.Id,
			Content:     message.Content,
			ContentType: int16(message.ContentType),
			Url:         message.Url,
		}
		db.Save(&saveMessage)
	}

}
