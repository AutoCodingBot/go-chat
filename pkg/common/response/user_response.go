package response

//用户消息返回
type UserResponse struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

//用户好友信息返回
type UserFriendResponse struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	// LatestMsg    string `json:"latestMsg"`
	OnlineStatus bool `json:"onlineStatus"`
}
