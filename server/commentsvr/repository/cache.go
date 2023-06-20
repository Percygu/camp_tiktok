package repository

import (
	"commentsvr/config"
	"commentsvr/constant"
	"commentsvr/middleware/cache"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

// SetCommentCacheInfo 给某一个video添加评论，评论以hash 形式存储
func SetCommentCacheInfo(comment *Comment) error {
	redisKey := constant.VideoInfoPrefix + strconv.FormatInt(comment.VideoId, 10)
	redisCli := cache.GetRedisCli()
	expired := time.Second * time.Duration(config.GetGlobalConfig().RedisConfig.Expired)

	commentIDStr := strconv.FormatInt(comment.VideoId, 10)
	err := redisCli.HSet(context.Background(), redisKey, commentIDStr, comment).Err()
	if err != nil {
		panic(err)
	}
	// 每有一个key的视频加入，就重新给key设置过期时间
	redisCli.Expire(context.Background(), redisKey, expired)
	return nil
}

// GetCommentCacheInfo 获取某个video的评论列表
func GetCommentCacheInfo(comment *Comment) error {
	redisKey := constant.CommentInfoPrefix + strconv.FormatInt(comment.CommentId, 10)

	val, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(config.GetGlobalConfig().RedisConfig.Expired)
	_, err = cache.GetRedisCli().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err
}
