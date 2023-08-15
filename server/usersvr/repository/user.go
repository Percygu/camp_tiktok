package repository

import (
	"errors"
)

// bool为true表示存在，false表示不存在
// 检查该用户名是否已经存在
func UserNameIsExist(userName string) (bool, error) {
	user, err := DbGetUserByUserName(userName)
	if err != nil {
		return false, err
	}
	if user.Id != 0 {
		return true, nil
	}
	return false, nil
}

// 创建用户
func InsertUser(userName, password string) (*User, error) {
	// 1. 插入数据库
	user, err := DbInsertUser(userName, password)
	if err != nil {
		return nil, err
	}

	// 2. 写入缓存
	go CacheSetUser(user)
	return &user, nil
}

func GetUserList(userIdList []int64) ([]*User, error) {
	users, err := DbGetUserList(userIdList)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 获取用户信息
func GetUserInfo(u interface{}) (user User, err error) {
	switch u := u.(type) {
	case int64:
		user, err = CacheGetUser(u)
		if err == nil {
			return user, nil
		}
		user, err = DbGetUserByUserId(u)
	case string:
		user, err = DbGetUserByUserName(u)
	default:
		err = errors.New("")
	}

	return user, err
}

func UpdateUserFavoritedNum(userID, updateType int64) error {
	var num int64
	// updateType 1: 点赞 else： 取消点赞
	if updateType == 1 {
		num = 1
	} else {
		num = -1
	}

	// 更新Db
	err := DbUpdateUserFavoritedNum(userID, num)
	if err != nil {
		return err
	}
	// 更新缓存（在原来的基础上，加上num）
	err = CacheIncrByUserFavoritedNum(userID, num)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserFavoriteNum(userID, updateType int64) error {
	var num int64
	// updateType 1: 点赞 else： 取消点赞
	if updateType == 1 {
		num = 1
	} else {
		num = -1
	}
	// 更新Db
	err := DbUpdateUserFavoriteNum(userID, num)
	if err != nil {
		return err
	}

	err = CacheIncrByUserFavoriteNum(userID, num)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserFollowNum(userID, updateType int64) error {
	var num int64
	// updateType 1: 关注 else： 取消关注
	if updateType == 1 {
		num = 1
	} else {
		num = -1
	}
	// 1. 更新Db
	err := DbUpdateUserFollowNum(userID, num)
	if err != nil {
		return err
	}

	// 2. 更新缓存（在原来的基础上，加上num）
	err = CacheIncrByUserFollowNum(userID, num)
	return nil
}

func UpdateUserFollowerNum(userID, updateType int64) error {
	var num int64
	// updateType 1: 关注 else： 取消关注
	if updateType == 1 {
		num = 1
	} else {
		num = -1
	}

	// 1. 更新Db
	err := DbUpdateUserFollowerNum(userID, num)
	if err != nil {
		return err
	}

	// 2. 更新缓存（在原来的基础上，加上num）
	err = CacheIncrByUserFollowerNum(userID, num)
	if err != nil {
		return err
	}
	return nil
}
