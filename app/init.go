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
	// register your controllers here
	r.GET("/api/getUserAccessToken", controller.GetUserAccessToken)
	r.POST("/api/project", controller.InitProject)

	//改为每个项目在自己的知识库单独储存会议，不需要持久化，弃用该项
	//r.POST("/api/meeting", controller.ReadMeetingForm)

	// DO NOT CHANGE LINES BELOW
	// register dispatcher
	r.POST("/feiShu/Event", dispatcher.Dispatcher)
}
