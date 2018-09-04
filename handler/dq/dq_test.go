package dq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"

	"github.com/TeamFat/DelayQueue/handler"
	"github.com/TeamFat/DelayQueue/pkg/redis"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	// 初始化路由
	router = gin.Default()
	router.Use(connRedis)
	router.POST("/queue/push", Push)
}

func connRedis(c *gin.Context) {
	viper.SetDefault("bucketKeyPrefix", "Test_DelayBucket_")
	viper.SetDefault("jobKeyPrefix", "Test_DelayJob_")
	viper.SetDefault("redis.redisInfo", "127.0.0.1:6379")
	viper.SetDefault("redis.redisAuth", 123456)
	viper.SetDefault("redis.redisDb", 2)
	viper.SetDefault("bucketCount", 4)
	redis.ConnRedis()
}

// PostJson 根据特定请求uri和参数param，以Json形式传递参数，发起post请求返回响应
func PostJson(uri string, param map[string]interface{}, router *gin.Engine) []byte {
	// 将参数转化为json比特流
	jsonByte, _ := json.Marshal(param)

	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))

	// 初始化响应
	w := httptest.NewRecorder()

	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	// 提取响应
	result := w.Result()
	defer result.Body.Close()

	// 读取响应body
	body, _ := ioutil.ReadAll(result.Body)
	return body
}

// TestPush 测试以Json形式访问入队的接口
func TestPush(t *testing.T) {
	// 初始化请求地址和请求参数
	uri := "/queue/push"

	param := make(map[string]interface{})
	param["topic"] = "order"
	param["delay"] = 30
	param["body"] = "{\"uid\": 10829378,\"created\": 1498657365 }"

	// 发起post请求，以Json形式传递参数
	body := PostJson(uri, param, router)
	fmt.Printf("response:%v\n", string(body))

	// 解析响应，判断响应是否与预期一致
	response := &handler.Response{}
	if err := json.Unmarshal(body, response); err != nil {
		t.Errorf("解析响应出错，err:%v\n", err)
	}
	if response.Code != 0 {
		t.Errorf("响应数据不符:%v\n", response.Message)
	}
}
