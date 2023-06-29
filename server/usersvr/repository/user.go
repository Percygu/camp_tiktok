package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"usersvr/middleware/cache"
	"usersvr/middleware/db"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 检查该用户名是否已经存在, 存在返回错误
func UserNameIsExist(userName string) (bool, error) {
	db := db.GetDB()
	user := User{}
	err := db.Where("user_name = ?", userName).First(&user).Error
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return false, err
		}
		return false, nil // 数据库错误
	}

	return true, nil
}

// 创建用户
func InsertUser(userName, password string) (*User, error) {
	db := db.GetDB()
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := User{
		Name:            userName,
		Password:        string(hashPassword),
		Follow:          0,
		Follower:        0,
		TotalFav:        0,
		FavCount:        0,
		Avatar:          "https://tse1-mm.cn.bing.net/th/id/R-C.d83ded12079fa9e407e9928b8f300802?rik=Gzu6EnSylX9f1Q&riu=http%3a%2f%2fwww.webcarpenter.com%2fpictures%2fGo-gopher-programming-language.jpg&ehk=giVQvdvQiENrabreHFM8x%2fyOU70l%2fy6FOa6RS3viJ24%3d&risl=&pid=ImgRaw&r=0",
		BackgroundImage: "https://tse2-mm.cn.bing.net/th/id/OIP-C.sDoybxmH4DIpvO33-wQEPgHaEq?pid=ImgDet&rs=1",
		Signature:       "test sign",
	}
	result := db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	zap.L().Info("create user", zap.Any("user", user))
	go CacheSetUser(user)
	return &user, nil
}

// 获取用户信息
func GetUserInfo(u interface{}) (User, error) {
	db := db.GetDB()
	user := User{}
	var err error
	switch u := u.(type) {
	case int64:
		user, err = CacheGetUser(u)
		if err == nil {
			return user, nil
		}
		err = db.Where("user_id = ?", u).Find(&user).Error

	case string:
		err = db.Where("user_name = ?", u).Find(&user).Error
	default:
		err = errors.New("")
	}
	if err != nil {
		return user, errors.New("user error")
	}

	go CacheSetUser(user)
	zap.L().Info("get user info", zap.Any("user", user))
	return user, nil
}

func CacheSetUser(u User) {
	uid := strconv.FormatInt(u.Id, 10)
	value, err := json.Marshal(u)
	if err != nil {
		zap.L().Error("json marshal error", zap.Error(err))
	}

	if cache.GetRedisCli().Set(context.Background(), "user_"+uid, value, cache.ValueExpire).Err(); err != nil {
		zap.L().Error("redis set error", zap.Error(err))
	}
}

func CacheGetUser(uid int64) (User, error) {
	key := strconv.FormatInt(uid, 10)

	data, err := cache.GetRedisCli().Get(context.Background(), "user_"+key).Bytes()
	user := User{}
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CacheHGet(key, mKey string) ([]byte, error) {

	data, err := cache.GetRedisCli().HGet(context.Background(), key, mKey).Bytes()
	if err != nil {
		return []byte{}, err
	}
	if len(data) == 0 {
		return []byte{}, errors.New("data is empty")
	}
	return data, nil
}
