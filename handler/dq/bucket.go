package dq

import "github.com/TeamFat/DelayQueue/pkg/redis"

// 添加JobId到bucket中
func pushToBucket(key string, timestamp int64, jobId string) error {
	_, err := redis.ExecRedisCommand("ZADD", key, timestamp, jobId)

	return err
}

// 从bucket中删除JobId
func removeFromBucket(bucket string, jobId string) error {
	_, err := redis.ExecRedisCommand("ZREM", bucket, jobId)

	return err
}
