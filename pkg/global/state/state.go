package state

import (
	"sort"
)

var (
	UserState = make([]string, 0)
)

func AppendOnline(uuid string) {
	if !noDuplicateInsurance(uuid) {
		return
	}
	UserState = append(UserState, uuid)
	sort.Strings(UserState)
}

func RemoveOnline(uuid string) {
	UserState = removeUUID(UserState, uuid)
}

// 确保新增数据不在已有的slice中
func noDuplicateInsurance(uuid string) bool {
	for _, val := range UserState {
		if val == uuid {
			return false
		}
	}
	return true
}

// removeUUID 从切片中移除指定的 uuid
func removeUUID(slice []string, uuid string) []string {
	var result []string
	for _, v := range slice {
		if v != uuid {
			result = append(result, v)
		}
	}
	sort.Strings(result)
	return result
}

// 某用户是否在线
func UserOnlineStatus(uuid string) bool {
	// 使用 sort.SearchStrings 进行二分查找
	i := sort.SearchStrings(UserState, uuid)

	return i < len(UserState) && UserState[i] == uuid
}
