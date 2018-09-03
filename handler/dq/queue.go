package dq

import (
	"github.com/TeamFat/DelayQueue/pkg/redis"
	"github.com/spf13/viper"
)

// 添加JobId到队列中
func pushToReadyQueue(queueName string, jobID string) error {
	queueName = viper.GetString("queueKeyPrefix") + queueName

	_, err := redis.ExecRedisCommand("RPUSH", queueName, jobID)

	return err
}

// 从队列中阻塞获取JobId
func blockPopFromReadyQueue(queues []string, timeout int) (string, error) {
	var args []interface{}
	for _, queue := range queues {
		queue = viper.GetString("queueKeyPrefix") + queue
		args = append(args, queue)
	}
	args = append(args, timeout)
	value, err := redis.ExecRedisCommand("BLPOP", args...)
	if err != nil {
		return "", err
	}
	if value == nil {
		return "", nil
	}
	var valueBytes []interface{}
	valueBytes = value.([]interface{})
	if len(valueBytes) == 0 {
		return "", nil
	}
	element := string(valueBytes[1].([]byte))

	return element, nil
}
