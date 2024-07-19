package response

import "time"

type MessageResponse struct {
	ID           int32     `json:"id" gorm:"primarykey"`
	FromUserId   int32     `json:"fromUserId" gorm:"index"`
	ToUserId     int32     `json:"toUserId" gorm:"index"`
	Content      string    `json:"content" gorm:"type:varchar(250)"`
	ContentType  int16     `json:"contentType" gorm:"comment:'消息内容类型：1文字，2语音，3视频'"`
	CreatedAt    time.Time `json:"createAt"`
	FromUsername string    `json:"fromUsername" gorm:"column:username"`
	ToUsername   string    `json:"toUsername"`
	Avatar       string    `json:"avatar"`
	Url          string    `json:"url"`
}
