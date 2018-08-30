package main

import (
	"errors"
	"flag"
	"net/http"
	"time"

	"github.com/TeamFat/DelayQueue/config"
	"github.com/TeamFat/DelayQueue/handler/dq"
	"github.com/TeamFat/DelayQueue/pkg/redis"
	"github.com/TeamFat/DelayQueue/router"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

var (
	cfg = flag.String("config", "", "apiserver config file path.")
)

func main() {
	flag.Parse()

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init redis pool
	if err := redis.ConnRedis(); err != nil {
		panic(err)
	}

	// init delay queue
	dq.Init()

	// Set gin mode.
	gin.SetMode(viper.GetString("runMode"))

	// Create the Gin engine.
	g := gin.New()

	middlewares := []gin.HandlerFunc{}

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlwares.
		middlewares...,
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			logs.Emergency("The router has no response, or it might took too long to start up.", err)
		}
		logs.Info("The router has been deployed successfully.")
	}()

	logs.Info("Start to listening the incoming requests on http address:", viper.GetString("addr"))
	logs.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("maxPingCount"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		logs.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
