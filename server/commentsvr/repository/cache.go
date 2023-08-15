package repository

import (
	"commentsvr/config"
	"commentsvr/log"
	"commentsvr/middleware/cache"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// CacheSetComment 给某一个video添加评论，评论以hash 形式存储
func CacheSetComment(comment *Comment) error {
	videoKey := VideoInfoPrefix + strconv.FormatInt(comment.VideoId, 10)
	redisCli := cache.GetRedisCli()
	// 过期时间
	expired := time.Second * time.Duration(config.GetGlobalConfig().RedisConfig.Expired)

	commentBytes, err := json.Marshal(comment)
	if err != nil {
		log.Errorf("CacheSetCommentInfo and json marshal comment err:%v", err)
		return err
	}
	commentIDStr := strconv.FormatInt(comment.VideoId, 10)
	err = redisCli.HSet(context.Background(), videoKey, commentIDStr, string(commentBytes)).Err()
	if err != nil {
		log.Error("CacheSetCommentInfo error", zap.Error(err))
		return err
	}
	// 每有一个key的视频加入，就重新给key设置过期时间
	redisCli.Expire(context.Background(), videoKey, expired)
	return nil
}

func CacheDelComment(keyList []string, VideoID int64) error {
	videoKey := VideoInfoPrefix + strconv.FormatInt(VideoID, 10)
	redisCli := cache.GetRedisCli()

	if err := redisCli.HDel(context.Background(), videoKey, keyList...).Err(); err != nil {
		log.Errorf("CacheDelComment error", zap.Error(err))
		return err
	}
	return nil
}

func CacheGetCommentList(vid int64) ([]*Comment, error) {
	videoKey := VideoInfoPrefix + strconv.FormatInt(vid, 10)
	redisCli := cache.GetRedisCli()
	resultList, err := redisCli.HGetAll(context.Background(), videoKey).Result()
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
