package app

import (
	"xlab-feishu-robot/app/controller"
	"xlab-feishu-robot/app/dispatcher"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	controller.InitEvent()
	Register(r)
}

func Register(r *gin.Engine) {
	r.GET("/api/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// DO NOT CHANGE LINES BELOW
	// register dispatcher
	r.POST("/feiShu/Event", dispatcher.Dispatcher)
}
