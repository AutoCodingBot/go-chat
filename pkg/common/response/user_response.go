package response

type UserResponse struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
