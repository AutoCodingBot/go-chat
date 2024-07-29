package userreq

type UserInfo struct {
	NickName        string `form:"nickname"`
	Avatar          string `form:"avatar"`
	CurrentPassword string `form:"currentPassword"`
	NewPassword     string `form:"newPassword"`
}
