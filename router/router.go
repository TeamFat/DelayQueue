package router

import (
	"github.com/TeamFat/DelayQueue/handler/dq"
	"github.com/TeamFat/DelayQueue/handler/sd"
	"github.com/TeamFat/DelayQueue/router/middleware"

	"github.com/TeamFat/DelayQueue/handler"
	"github.com/TeamFat/DelayQueue/pkg/errno"
	"github.com/gin-gonic/gin"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	//g.Use(middleware.NoCache)
	//g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		//c.String(http.StatusNotFound, "The incorrect API route.")
		handler.SendResponse(c, errno.StatusNotFound, nil)
		return
	})

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	dqrouter := g.Group("/queue")
	{
		dqrouter.POST("/push", dq.Push)
		dqrouter.POST("/pop", dq.Pop)
	}

	return g
}
