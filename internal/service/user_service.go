package service

import (
	"time"

	"chat-room/internal/dao/pool"
	"chat-room/internal/model"
	"chat-room/pkg/common/request"
	"chat-room/pkg/common/response"
	"chat-room/pkg/errors"
	"chat-room/pkg/global/log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
}

var UserService = new(userService)

func (u *userService) Register(user *model.User) error {
	db := pool.GetDB()
	var userCount int64
	db.Model(user).Where("username", user.Username).Where("email", user.Email).Where("nickname", user.Nickname).Count(&userCount)
	if userCount > 0 {
		return errors.New("Duplicate Email Or UserName or NickName")
	}
	user.Uuid = uuid.New().String()

	//encrype Password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return errors.New("Process failed")
	}
	user.Password = hashedPassword

	user.CreateAt = time.Now()
	user.DeleteAt = 0

	db.Create(&user)
	return nil
}

func (u *userService) Login(user *model.User) error {
	pool.GetDB().AutoMigrate(&user)

	db := pool.GetDB()

	var queryUser *model.User
	db.First(&queryUser, "username = ?", user.Username)
	log.Logger.Debug("queryUSer", log.Any("user in service", queryUser))
	if queryUser.Id == 0 {
		return errors.New("Invalid User OR Password")
	}
	//validate password
	res := checkPasswordHash(user.Password, queryUser.Password)
	if !res {
		return errors.New("Invalid User OR Password")
	}
	user.Uuid = queryUser.Uuid
	user.Id = queryUser.Id
	return nil
}

func (u *userService) ModifyUserInfo(user *model.User) error {
	var queryUser *model.User
	db := pool.GetDB()
	db.First(&queryUser, "username = ?", user.Username)
	log.Logger.Debug("queryUser", log.Any("queryUser", queryUser))
	var nullId int32 = 0
	if nullId == queryUser.Id {
		return errors.New("用户不存在")
	}
	queryUser.Nickname = user.Nickname
	queryUser.Email = user.Email
	queryUser.Password = user.Password

	db.Save(queryUser)
	return nil
}

func (u *userService) GetUserDetails(uuid string) model.User {
	var queryUser *model.User
	db := pool.GetDB()
	db.Select("uuid", "username", "nickname", "avatar").First(&queryUser, "uuid = ?", uuid)
	return *queryUser
}

// 通过名称查找群组或者用户
func (u *userService) GetUserOrGroupByName(name string) response.SearchResponse {
	var queryUser *model.User
	db := pool.GetDB()
	db.Select("uuid", "username", "nickname", "avatar").First(&queryUser, "username = ?", name)

	var queryGroup *model.Group
	db.Select("uuid", "name").First(&queryGroup, "name = ?", name)

	search := response.SearchResponse{
		User:  *queryUser,
		Group: *queryGroup,
	}
	return search
}

func (u *userService) GetUserList(uuid string) []model.User {
	db := pool.GetDB()

	var queryUser *model.User
	db.First(&queryUser, "uuid = ?", uuid)
	var nullId int32 = 0
	if nullId == queryUser.Id {
		return nil
	}

	var queryUsers []model.User
	db.Raw("SELECT u.username, u.uuid, u.avatar FROM user_friends AS uf JOIN users AS u ON uf.friend_id = u.id WHERE uf.user_id = ?", queryUser.Id).Scan(&queryUsers)

	return queryUsers
}

func (u *userService) AddFriend(userFriendRequest *request.FriendRequest) error {
	var queryUser *model.User
	db := pool.GetDB()
	db.First(&queryUser, "uuid = ?", userFriendRequest.Uuid)
	log.Logger.Debug("queryUser", log.Any("queryUser", queryUser))
	var nullId int32 = 0
	if nullId == queryUser.Id {
		return errors.New("用户不存在")
	}

	var friend *model.User
	db.First(&friend, "username = ?", userFriendRequest.FriendUsername)
	if nullId == friend.Id {
		return errors.New("已添加该好友")
	}

	userFriend := model.UserFriend{
		UserId:   queryUser.Id,
		FriendId: friend.Id,
	}

	var userFriendQuery *model.UserFriend
	db.First(&userFriendQuery, "user_id = ? and friend_id = ?", queryUser.Id, friend.Id)
	if userFriendQuery.ID != nullId {
		return errors.New("该用户已经是你好友")
	}

	db.AutoMigrate(&userFriend)
	db.Save(&userFriend)
	log.Logger.Debug("userFriend", log.Any("userFriend", userFriend))

	return nil
}

// 修改头像
func (u *userService) ModifyUserAvatar(avatar string, uid int, objectType string) error {
	db := pool.GetDB()
	if objectType == "user" {
		var targetModel model.User
		db.Debug().Model(&targetModel).Where("id=?", uid).Update("avatar", avatar)

	} else if objectType == "group" {
		// log.Logger.Debug("2", log.Any("2-1", 2))

		var targetModel model.Group
		// Is the owner of group?
		var exiests int64
		db.Model(&targetModel).Where("user_id = ?", uid).Count(&exiests)
		if exiests < 1 {
			return errors.New("U are not the Owner of this group")
		}
		db.Model(&targetModel).Update("avatar", avatar)

	} else {
		return errors.New("未选择对象")
	}

	return nil
}

// 使用bcrypt哈希运算密码
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 验证密码
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
