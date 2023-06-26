package repository

import (
	"commentsvr/config"
	"commentsvr/constant"
	"commentsvr/log"
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

	commentBytes, err := json.Marshal(comment)
	if err != nil {
		log.Errorf("json marshal comment err:%v", err)
		return err
	}
	commentIDStr := strconv.FormatInt(comment.VideoId, 10)
	err = redisCli.HSet(context.Background(), redisKey, commentIDStr, string(commentBytes)).Err()
	if err != nil {
		log.Errorf("redis hset comment err:%v", err)
		return err
	}
	// 每有一个key的视频加入，就重新给key设置过期时间
	redisCli.Expire(context.Background(), redisKey, expired)
	return nil
}

func DelCommentCacheInfo(keyList []string, VideoID int64) error {
	redisKey := constant.VideoInfoPrefix + strconv.FormatInt(VideoID, 10)
	redisCli := cache.GetRedisCli()
	if err := redisCli.HDel(context.Background(), redisKey, keyList...).Err(); err != nil {
		log.Errorf("del redis hkey err:%v", err)
		return err
	}
	return nil
}

// GetCommentCacheInfo 获取某个video的评论列表
func GetCommentCacheInfo(comment *Comment) error {
	redisKey := constant.CommentInfoPrefix + strconv.FormatInt(comment.Id, 10)

	val, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(config.GetGlobalConfig().RedisConfig.Expired)
	_, err = cache.GetRedisCli().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err
}

func GetCommentCacheList(vid int64) ([]*Comment, error) {
	// key名 video:1 value: comment1 comment2 comment3
	redisKey := constant.VideoInfoPrefix + strconv.FormatInt(vid, 10)
	redisCli := cache.GetRedisCli()
	resultList, err := redisCli.HGetAll(context.Background(), redisKey).Result()
	if err != nil {
		log.Errorf("HGet all video with %d comments err:%v", vid, err)
		return nil, err
	}
	var commentList []*Comment
	for _, v := range resultList {
		var comment Comment
		if err := json.Unmarshal([]byte(v), &comment); err != nil {
			log.Errorf("json unmarshal comment with %s err:%v", v, err)
			return nil, err
		}
		commentList = append(commentList, &comment)
	}
	return commentList, nil
}

func DelCacheCommentAll(vid int64) error {
	redisKey := constant.VideoInfoPrefix + strconv.FormatInt(vid, 10)
	redisCli := cache.GetRedisCli()
	if err := redisCli.Del(context.Background(), redisKey).Err(); err != nil {
		log.Errorf("del video with %d all comments err:%v", vid, err)
		return err
	}
	return nil
}
