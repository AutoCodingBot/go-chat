package server

import (
	"chat-room/config"
	"chat-room/internal/kafka"
	"chat-room/pkg/common/constant"
	"chat-room/pkg/global/log"
	"chat-room/pkg/protocol"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Name string
	Send chan []byte
}

func (c *Client) Read() {
	defer func() {
		MyServer.Ungister <- c
		c.Conn.Close()
	}()

	for {
		c.Conn.PongHandler()
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Logger.Error("client read message error", log.Any("client read message error", err.Error()))
			MyServer.Ungister <- c
			c.Conn.Close()
			break
		}

		msg := &protocol.Message{}
		proto.Unmarshal(message, msg)

		// pong
		if msg.Type == constant.HEAT_BEAT {
			pong := &protocol.Message{
				Content: constant.PONG,
				Type:    constant.HEAT_BEAT,
			}
			pongByte, err2 := proto.Marshal(pong)
			if nil != err2 {
				log.Logger.Error("client marshal message error", log.Any("client marshal message error", err2.Error()))
			}
			c.Conn.WriteMessage(websocket.BinaryMessage, pongByte)
		} else {
			if config.GetConfig().MsgChannelType.ChannelType == constant.KAFKA {
				kafka.Send(message)
			} else {
				MyServer.Broadcast <- message
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		c.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}

/*
这段代码定义了一个`Client`结构体以及其关联的`Read`和`Write`方法，用于处理WebSocket连接上的数据读写逻辑，服务于一个聊天室应用。以下是详细解析：

### `Client` 结构体

- `Conn`: 表示与客户端相连的WebSocket连接。
- `Name`: 存储客户端的名称，用于识别用户。
- `Send`: 一个channel，用于从服务器向客户端发送消息。

### `Read` 方法

此方法负责从WebSocket连接中读取消息，并处理这些消息。它的工作流程如下：

1. **关闭清理**：使用`defer`语句确保在方法结束时，无论是否发生错误，都能正确地从服务器注销客户端，并关闭WebSocket连接。
2. **心跳处理**：通过调用`c.Conn.PongHandler()`设置自动响应Ping帧，维护连接活跃状态。随后读取消息，如果读取失败，则执行清理操作并退出循环。
3. **消息解码与处理**：将接收到的二进制消息解码为`protocol.Message`结构体。对于心跳（`HEAT_BEAT`）消息，构造并发送Pong响应。其他消息根据配置决定处理方式：若配置为使用Kafka，则将消息发送到Kafka；否则，将消息放入服务器的广播channel中。

### `Write` 方法

此方法在独立的goroutine中运行，负责从`Send` channel中读取消息并将其写回到WebSocket连接，直至channel被关闭。

### 关键技术点

- **WebSocket**: 使用标准的WebSocket协议为客户端提供全双工通信，适用于实时聊天应用。
- **goroutine**: 通过在`Read`和`Write`方法中启动独立的goroutine，实现并发读写，提高服务器的响应能力。
- **心跳机制**：通过设置Pong处理器保持WebSocket连接活跃，防止因长时间无数据交换导致的连接中断。
- **消息中间件集成**：根据配置，可以选择将消息直接广播或通过Kafka消息队列处理，增加了系统的可扩展性和灵活性。
- **protobuf**: 利用Protocol Buffers进行消息序列化和反序列化，提高了消息传输的效率和兼容性。

整体而言，这段代码展示了聊天室服务器如何高效地处理客户端的WebSocket连接，包括消息的接收、心跳维护、以及消息的分发或异步处理，体现了现代实时通讯服务的关键技术实践。
*/
