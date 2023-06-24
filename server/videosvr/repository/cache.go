package repository

import (
	"context"
	"strconv"
	"videosvr/middleware/cache"
)

func CacheSetAuthor(videoId, authorId int64) error {
	err := cache.GetRedisCli().HSet(context.Background(), "video", strconv.FormatInt(videoId, 10), authorId).Err()
	if err != nil {
		return err
	}
	return nil
}
