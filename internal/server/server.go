package server

import (
	"chat-room/config"
	"chat-room/internal/service"
	"chat-room/pkg/common/constant"
	"chat-room/pkg/common/util"
	"chat-room/pkg/global/log"
	"chat-room/pkg/global/state"
	"chat-room/pkg/protocol"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
)

var MyServer = NewServer()

type Server struct {
	Clients   map[string]*Client
	mutex     *sync.Mutex
	Broadcast chan []byte
	Register  chan *Client
	Ungister  chan *Client
}

func NewServer() *Server {
	return &Server{
		mutex:     &sync.Mutex{},
		Clients:   make(map[string]*Client),
		Broadcast: make(chan []byte),
		Register:  make(chan *Client),
		Ungister:  make(chan *Client),
	}
}

// 消费kafka里面的消息, 然后直接放入go channel中统一进行消费
func ConsumerKafkaMsg(data []byte) {
	MyServer.Broadcast <- data
}

func (s *Server) Start() {
	log.Logger.Info("start server", log.Any("start server", "start server..."))
	for {
		select {
		case conn := <-s.Register:
			log.Logger.Info("login", log.Any("login", "new user login in"+conn.Name))
			state.AppendOnline(conn.Name)
			s.Clients[conn.Name] = conn
			msg := &protocol.Message{
				From:    "System",
				To:      conn.Name,
				Content: "welcome!",
			}
			protoMsg, _ := proto.Marshal(msg)
			conn.Send <- protoMsg

		case conn := <-s.Ungister:
			log.Logger.Info("loginout", log.Any("loginout", conn.Name))
			state.RemoveOnline(conn.Name)
			if _, ok := s.Clients[conn.Name]; ok {
				close(conn.Send)
				delete(s.Clients, conn.Name)
			}

		case message := <-s.Broadcast:
			msg := &protocol.Message{}
			proto.Unmarshal(message, msg)

			if msg.To != "" {
				// 一般消息，比如文本消息，视频文件消息等
				if msg.ContentType >= constant.TEXT && msg.ContentType <= constant.VIDEO {
					// 保存消息只会在存在socket的一个端上进行保存，防止分布式部署后，消息重复问题
					_, exits := s.Clients[msg.From]
					if exits {
						saveMessage(msg)
					}

					if msg.MessageType == constant.MESSAGE_TYPE_USER {
						client, ok := s.Clients[msg.To]
						if ok {
							msgByte, err := proto.Marshal(msg)
							if err == nil {
								client.Send <- msgByte
							}
						}
					} else if msg.MessageType == constant.MESSAGE_TYPE_GROUP {
						sendGroupMessage(msg, s)
					}
				} else {
					// 语音电话，视频电话等，仅支持单人聊天，不支持群聊
					// 不保存文件，直接进行转发
					client, ok := s.Clients[msg.To]
					if ok {
						client.Send <- message
					}
				}

			} else {
				// 无对应接受人员进行广播
				for id, conn := range s.Clients {
					log.Logger.Info("allUser", log.Any("allUser", id))

					select {
					case conn.Send <- message:
					default:
						close(conn.Send)
						delete(s.Clients, conn.Name)
					}
				}
			}
		}
	}
}

// 发送给群组消息,需要查询该群所有人员依次发送
func sendGroupMessage(msg *protocol.Message, s *Server) {
	// 发送给群组的消息，查找该群所有的用户进行发送
	users := service.GroupService.GetUserIdByGroupUuid(msg.To)
	for _, user := range users {
		if user.Uuid == msg.From {
			continue
		}

		client, ok := s.Clients[user.Uuid]
		if !ok {
			continue
		}

		fromUserDetails := service.UserService.GetUserDetails(msg.From)
		// 由于发送群聊时，from是个人，to是群聊uuid。所以在返回消息时，将form修改为群聊uuid，和单聊进行统一
		msgSend := protocol.Message{
			Avatar:       fromUserDetails.Avatar,
			FromUsername: msg.FromUsername,
			From:         msg.To,
			To:           msg.From,
			Content:      msg.Content,
			ContentType:  msg.ContentType,
			Type:         msg.Type,
			MessageType:  msg.MessageType,
			Url:          msg.Url,
		}

		msgByte, err := proto.Marshal(&msgSend)
		if err == nil {
			client.Send <- msgByte
		}
	}
}

// 保存消息，如果是文本消息直接保存，如果是文件，语音等消息，保存文件后，保存对应的文件路径
func saveMessage(message *protocol.Message) {
	// 如果上传的是base64字符串文件，解析文件保存
	if message.ContentType == 2 {
		url := uuid.New().String() + ".png"
		index := strings.Index(message.Content, "base64")
		index += 7

		content := message.Content
		content = content[index:]

		dataBuffer, dataErr := base64.StdEncoding.DecodeString(content)
		if dataErr != nil {
			log.Logger.Error("transfer base64 to file error", log.String("transfer base64 to file error", dataErr.Error()))
			return
		}
		err := ioutil.WriteFile(config.GetConfig().StaticPath.FilePath+url, dataBuffer, 0666)
		if err != nil {
			log.Logger.Error("write file error", log.String("write file error", err.Error()))
			return
		}
		message.Url = url
		message.Content = ""
	} else if message.ContentType == 3 {
		// 普通的文件二进制上传
		fileSuffix := util.GetFileType(message.File)
		nullStr := ""
		if nullStr == fileSuffix {
			fileSuffix = strings.ToLower(message.FileSuffix)
		}
		contentType := util.GetContentTypeBySuffix(fileSuffix)
		url := uuid.New().String() + "." + fileSuffix
		err := ioutil.WriteFile(config.GetConfig().StaticPath.FilePath+url, message.File, 0666)
		if err != nil {
			log.Logger.Error("write file error", log.String("write file error", err.Error()))
			return
		}
		message.Url = url
		message.File = nil
		message.ContentType = contentType
	}

	service.MessageService.SaveMessage(*message)
}

/*
这段代码是一个基于Go语言实现的聊天室服务器的核心逻辑部分，采用了gRPC的protobuf协议进行消息编码与解码，利用了goroutine和channel来处理并发和消息传递。以下是关键功能的解析：

### 主要结构和变量

- `Server` 结构体定义了聊天室服务器的主要属性，包括：
  - `Clients`: 一个映射，存储所有已连接客户端（基于用户名为键）。
  - `mutex`: 用于保护`Clients`映射的互斥锁。
  - `Broadcast`, `Register`, `Unregister`: 三个channel，分别用于广播消息、注册新客户端、注销客户端。

- `NewServer` 函数用于初始化`Server`实例。

### 关键函数

- `ConsumerKafkaMsg` 函数模拟了从Kafka等消息队列中消费消息，并将消息直接放入`Broadcast` channel中，供服务器处理。

- `Start` 方法是服务器的主循环，它持续监听三个channel，处理客户端的注册、注销、以及消息广播：
  - 当有新客户端注册时，将其添加至`Clients`映射，并发送欢迎消息。
  - 当客户端注销时，关闭其发送通道并从映射中移除。
  - 对于接收到的消息，根据消息的类型和目标（单聊或群聊），进行相应的处理和转发。

- `sendGroupMessage` 函数处理群聊消息，通过查询数据库获取群组内所有成员，然后逐一向他们发送消息。

- `saveMessage` 函数负责保存消息到数据库，特别是处理文件消息（如图片、普通文件等），先将文件保存到磁盘，再更新消息内容为文件的URL，并最终保存到数据库。

### 技术要点

- **并发处理**：通过goroutine和channel实现高并发消息处理，保证了系统的可扩展性和响应速度。
- **协议缓冲**：使用protobuf进行消息序列化和反序列化，提高效率并简化跨语言通信。
- **消息分发**：根据消息目标（单播或广播）灵活处理，实现了聊天室的基础通信逻辑。
- **数据持久化**：通过调用服务层函数（如`service.MessageService.SaveMessage`）保存聊天记录到数据库，保证聊天历史的可追溯性。
- **文件处理**：针对不同类型的文件消息（如图片、文档）提供了基础的文件上传和存储逻辑，包括文件类型检测和保存路径管理。

综上，这段代码展示了构建一个基础聊天室服务器的核心逻辑，涉及消息传递、用户管理、群聊支持及文件处理等核心功能。
*/
