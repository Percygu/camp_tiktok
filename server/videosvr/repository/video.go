package repository

import (
	"gorm.io/gorm"
)

func GetVideoList(AuthorId int64) ([]Video, error) {
	var videos []Video

	// TODO:
	// author, err := repository.GetUserInfo(AuthorId)
	if err != nil {
		return videos, err
	}
	db := global.DB
	err = db.Where("author_id = ?", AuthorId).Order("video_id DESC").Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return videos, err
	}
	for i := range videos {
		videos[i].Author = author
	}
	return videos, nil
}
