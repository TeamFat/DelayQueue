package dq

import (
	"github.com/TeamFat/DelayQueue/pkg/redis"
	"github.com/vmihailenco/msgpack"
)

// Job 使用msgpack序列化后保存到Redis,减少内存占用
type Job struct {
	Topic string `json:"topic" msgpack:"1"`
	ID    string `json:"id" msgpack:"2"`    // job唯一标识ID
	Delay int64  `json:"delay" msgpack:"3"` // 延迟时间, unix时间戳
	Body  string `json:"body" msgpack:"4"`
}

// 获取Job
func getJob(key string) (*Job, error) {
	value, err := redis.ExecRedisCommand("GET", key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	byteValue := value.([]byte)
	job := &Job{}
	err = msgpack.Unmarshal(byteValue, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

// 添加Job
func putJob(key string, job *Job) error {
	value, err := msgpack.Marshal(job)
	if err != nil {
		return err
	}
	_, err = redis.ExecRedisCommand("SET", key, value)

	return err
}

// 删除Job
func removeJob(key string) error {
	_, err := redis.ExecRedisCommand("DEL", key)

	return err
}
