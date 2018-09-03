package dq

import (
	"fmt"
	"hash/crc32"
	"strings"
	"time"

	"github.com/TeamFat/DelayQueue/handler"
	"github.com/TeamFat/DelayQueue/pkg/errno"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

type jobPush struct {
	Topic string `json:"topic" binding:"required"`
	Delay int64  `json:"delay" binding:"required"` //秒
	Body  string `json:"body" binding:"required"`
}

type jobPop struct {
	Topic   string `json:"topic" binding:"required"`
	Timeout int    `json:"timeout"`
}

// 定时器数量和bucket数量相同一一对应
var timers []*time.Ticker

// Push 入队
func Push(c *gin.Context) {
	var jobP jobPush

	if err := c.ShouldBindJSON(&jobP); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}
	var job Job
	job.Topic = strings.TrimSpace(jobP.Topic)
	job.Body = strings.TrimSpace(jobP.Body)
	if job.Topic == "" {
		handler.SendResponse(c, errno.ErrValidationTopic, nil)
		return
	}
	if jobP.Delay <= 0 || jobP.Delay > (1<<31) {
		handler.SendResponse(c, errno.ErrValidationDelay, nil)
		return
	}
	if job.Body == "" {
		handler.SendResponse(c, errno.ErrValidationBody, nil)
		return
	}
	job.Delay = time.Now().Unix() + jobP.Delay
	u, err := uuid.NewV4()
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, err)
		logs.Error(err)
		return
	}
	job.ID = viper.GetString("jobKeyPrefix") + u.String()
	logs.Info(job)
	err = putJob(job.ID, job)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, err)
		logs.Error(err)
		return
	}
	crc32q := crc32.MakeTable(0xD5828281)
	i := crc32.Checksum([]byte(job.Body), crc32q) % uint32(viper.GetInt("bucketCount"))
	bucketName := fmt.Sprintf(viper.GetString("bucketKeyPrefix")+"%d", i+1)
	err = pushToBucket(bucketName, job.Delay, job.ID)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, err)
		logs.Error(err)
		return
	}
	handler.SendResponse(c, errno.OK, nil)
	return
}

// Pop 出队
func Pop(c *gin.Context) {
	var jobP jobPop

	if err := c.ShouldBindJSON(&jobP); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}
	logs.Info(jobP)
	if jobP.Timeout == 0 {
		jobP.Timeout = viper.GetInt("queueBlockTimeout")
	}
	logs.Info(jobP)
	topics := strings.Split(jobP.Topic, ",")
	jobID, err := blockPopFromReadyQueue(topics, jobP.Timeout)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, err)
		logs.Error(err)
		return
	}
	if jobID == "" {
		handler.SendResponse(c, errno.OK, nil)
		return
	}
	job, err := getJob(jobID)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, err)
		logs.Error(err)
		return
	}
	if job == nil {
		handler.SendResponse(c, errno.OK, nil)
		return
	}
	err = removeJob(job.ID)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, err)
		logs.Error(err)
		return
	}
	handler.SendResponse(c, errno.OK, job)
	return
}

// Init 延时队列初始化
func Init() {
	//初始化定时器
	initTimers()
}

func initTimers() {
	timers = make([]*time.Ticker, viper.GetInt("bucketCount"))
	var bucketName string
	for i := 0; i < viper.GetInt("bucketCount"); i++ {
		timers[i] = time.NewTicker(1 * time.Second)
		bucketName = fmt.Sprintf(viper.GetString("bucketKeyPrefix")+"%d", i+1)
		go waitTicker(timers[i], bucketName)
	}
}

func waitTicker(timer *time.Ticker, bucketName string) {
	for {
		select {
		case t := <-timer.C:
			tickHandler(t, bucketName)
		}
	}
}

// 扫描bucket, 取出延迟时间小于当前时间的Job
func tickHandler(t time.Time, bucketName string) {
	for {
		bucketItem, err := getFromBucket(bucketName)
		if err != nil {
			logs.Error("扫描bucket错误", bucketName, err.Error())
			return
		}

		// 集合为空
		if bucketItem == nil {
			return
		}

		// 延迟时间未到
		if bucketItem.timestamp > t.Unix() {
			return
		}

		// 延迟时间小于等于当前时间, 取出Job元信息并放入ready queue
		job, err := getJob(bucketItem.jobID)
		if err != nil {
			logs.Error("获取Job元信息失败", bucketName, err.Error())
			continue
		}

		// job元信息不存在, 从bucket中删除
		if job == nil {
			removeFromBucket(bucketName, bucketItem.jobID)
			continue
		}

		err = pushToReadyQueue(job.Topic, bucketItem.jobID)
		if err != nil {
			logs.Error("JobId放入ready queue失败", bucketName, job, err.Error())
			continue
		}

		// 从bucket中删除
		removeFromBucket(bucketName, bucketItem.jobID)
	}
}
