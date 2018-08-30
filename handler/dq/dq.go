package dq

import (
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 定时器数量和bucket数量相同一一对应
var timers []*time.Ticker

// Push 入队
func Push(c *gin.Context) {
	message := "Push"
	c.String(http.StatusOK, message)
}

// Pop 出队
func Pop(c *gin.Context) {
	message := "Push"
	c.String(http.StatusOK, message)
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
	logs.Info(t, bucketName)
}
