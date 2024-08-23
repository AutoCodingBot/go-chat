package service

import (
	"time"

	"chat-room/internal/dao/pool"
	"chat-room/internal/model"
	userreq "chat-room/internal/request"
	"chat-room/internal/utils"
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
	// log.Logger.Debug("queryUSer", log.Any("user in service", queryUser))
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

func (u *userService) GetUserDetails(uuid string) response.UserResponse {
	var moselUser *model.User
	var responseUser *response.UserResponse
	db := pool.GetDB()
	db.Model(moselUser).Select("uuid", "username", "nickname", "avatar").Where("uuid = ?", uuid).First(&responseUser)
	return *responseUser
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

func (u *userService) GetUserList(claims utils.JwtCustClaim) []response.UserResponse {
	db := pool.GetDB()

	var users []model.User
	var modelFriends []model.UserFriend
	var response []response.UserResponse
	//get all friends id
	var activeAddition []int   //发起添加
	var positiveAddition []int //被添加
	db.Model(modelFriends).Select("friend_id").Where("user_id =?", claims.ID).Scan(&activeAddition)
	db.Model(modelFriends).Select("user_id").Where("friend_id =?", claims.ID).Scan(&positiveAddition)
	uidList := append(activeAddition, positiveAddition...)
	//fetch friendsInfo
	db.Model(users).Select("username,uuid,avatar").Where("id IN (?)", uidList).Scan(&response)
	return response

	// log.Logger.Debug("2", log.Any("2-1", activeAddition))
	// log.Logger.Debug("2", log.Any("2-2", positiveAddition))
	db.Debug().Raw(`
        WITH friends AS (
            SELECT DISTINCT uf.friend_id AS id
            FROM user_friends uf
            WHERE uf.user_id = ? OR uf.friend_id = ?
        )
        SELECT u.username, u.uuid, u.avatar
        FROM users u
        JOIN friends ON u.id = friends.id
    `, claims.ID, claims.ID).Scan(&response)
	/*
			SELECT
			u.username,
			u.uuid,
			u.avatar
		FROM
			( SELECT DISTINCT uf.friend_id AS id FROM user_friends uf WHERE uf.user_id = 6 OR uf.friend_id = 6 ) AS friends
			JOIN users u ON u.id = friends.id;
	*/
	return response
}

func (u *userService) AddFriend(userInfo utils.JwtCustClaim, userFriendRequest *request.FriendRequest) (model.User, error) {
	db := pool.GetDB()
	//freindInfo
	var friendInfo *model.User
	db.First(&friendInfo, "username = ?", userFriendRequest.FriendUsername)
	if friendInfo.Id == 0 {
		return model.User{}, errors.New("User did not exists!")
	}
	//Added already?
	var friendRelation model.UserFriend
	var friendExist int64
	relationId := utils.GenerateConversationId(userInfo.UserName, userFriendRequest.FriendUsername)
	db.Model(friendRelation).Where("relation_id = ?", relationId).Count(&friendExist)
	if friendExist > 0 {
		return model.User{}, errors.New("Friend added already!")
	}
	// add friend
	insertData := model.UserFriend{UserId: int32(userInfo.ID), FriendId: friendInfo.Id, RelationId: relationId}
	result := db.Create(&insertData) // 通过数据的指针来创建
	if result.Error != nil {
		return model.User{}, errors.New("Something goes wrong,try later")
	}
	return *friendInfo, nil
	// insertData.ID             // 返回插入数据的主键
	// result.Error        // 返回 error
	// result.RowsAffected // 返回插入记录的条数
}

// 修改头像
// func (u *userService) ModifyUserAvatar(avatar string, uid int, objectType string) error {
// 	db := pool.GetDB()
// 	if objectType == "user" {
// 		var targetModel model.User
// 		db.Debug().Model(&targetModel).Where("id=?", uid).Update("avatar", avatar)

// 	} else if objectType == "group" {
// 		// log.Logger.Debug("2", log.Any("2-1", 2))

// 		var targetModel model.Group
// 		// Is the owner of group?
// 		var exiests int64
// 		db.Model(&targetModel).Where("user_id = ?", uid).Count(&exiests)
// 		if exiests < 1 {
// 			return errors.New("U are not the Owner of this group")
// 		}
// 		db.Model(&targetModel).Update("avatar", avatar)

// 	} else {
// 		return errors.New("未选择对象")
// 	}

//		return nil
//	}
func (u *userService) UpdateUserProfile(claims utils.JwtCustClaim, userInfo *userreq.UserInfo) error {
	db := pool.GetDB()

	var queryUser *model.User
	db.First(&queryUser, "id = ?", claims.ID)

	//edit update data
	encryptedPassword, err := hashPassword(userInfo.NewPassword)
	if err != nil {
		return errors.New("Something goes wrong")
	}
	//update
	var updateData model.User
	//update with password
	if userInfo.NewPassword != "" {
		//validate password
		passwprdCheck := checkPasswordHash(userInfo.CurrentPassword, queryUser.Password)
		if !passwprdCheck {
			return errors.New("Incorrect Password")
		}
		updateData.Avatar = userInfo.Avatar
		updateData.Password = encryptedPassword
		//update without password
	} else {
		updateData.Avatar = userInfo.Avatar
	}
	db.Model(&queryUser).Where("id = ?", claims.ID).Updates(updateData)

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
