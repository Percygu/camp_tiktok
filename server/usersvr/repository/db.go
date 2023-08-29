package repository

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"usersvr/middleware/db"
)

func DbGetUserByUserId(userId int64) (User, error) {
	db := db.GetDB()
	user := User{}
	err := db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		zap.L().Error("DbGetUserByUserId error", zap.Error(err))
		return User{}, err
	}
	return user, nil
}

func DbGetUserByUserName(userName string) (User, error) {
	db := db.GetDB()
	user := User{}
	err := db.Where("user_name = ?", userName).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		zap.L().Error("DbGetUserByUserName error", zap.Error(err))
		return User{}, err
	}
	return user, nil
}

func DbInsertUser(userName, password string) (User, error) {
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
	result := db.Model(&User{}).Create(&user)
	if result.Error != nil {
		zap.L().Error("DbInsertUser error", zap.Error(result.Error))
		return User{}, result.Error
	}
	return user, nil
}

func DbUpdateUserFavoritedNum(userId int64, num int64) error {
	db := db.GetDB()
	err := db.Model(&User{}).Where("id = ?", userId).Update("total_favorited", gorm.Expr("total_favorited + ?", num)).Error
	if err != nil {
		zap.L().Error("DbUpdateUserFavoritedNum error", zap.Error(err))
		return err
	}
	return nil
}

func DbUpdateUserFavoriteNum(userId int64, num int64) error {
	db := db.GetDB()
	err := db.Model(&User{}).Where("id = ?", userId).Update("favorite_count", gorm.Expr("favorite_count + ?", num)).Error
	if err != nil {
		zap.L().Error("DbUpdateUserFavoriteNum error", zap.Error(err))
		return err
	}
	return nil
}

func DbUpdateUserFollowNum(userId int64, num int64) error {
	db := db.GetDB()
	err := db.Model(&User{}).Where("id = ?", userId).Update("follow_count", gorm.Expr("follow_count + ?", num)).Error
	if err != nil {
		zap.L().Error("DbUpdateUserFollowNum error", zap.Error(err))
		return err
	}
	return nil
}

func DbGetUserList(userIdList []int64) ([]*User, error) {
	db := db.GetDB()
	var users []*User
	err := db.Where("id in ?", userIdList).Find(&users).Error
	if err != nil {
		zap.L().Error("DbGetUserList error", zap.Error(err))
		return nil, err
	}
	return users, nil

}

func DbUpdateUserFollowerNum(userId int64, num int64) error {
	db := db.GetDB()
	err := db.Model(&User{}).Where("id = ?", userId).Update("follower_count", gorm.Expr("follower_count + ?", num)).Error
	if err != nil {
		zap.L().Error("DbUpdateUserFollowerNum error", zap.Error(err))
		return err
	}
	return nil
}
